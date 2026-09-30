require "cgi"

module ApplicationHelper
  # The chapter bodies carry intentional inline HTML (lists, <em>, <aside>
  # subtitles), so it is rendered and then sanitized against an allowlist instead
  # of being escaped or trusted blindly.
  MARKDOWN_TAGS = %w[
    p br hr ul ol li em strong i b blockquote cite code pre
    h1 h2 h3 h4 h5 h6 a sup sub aside div span mark small
    table thead tbody tr td th dl dt dd
  ].freeze

  MARKDOWN_ATTRIBUTES = {
    "a" => %w[href title rel target],
    "aside" => %w[class],
    "div" => %w[class],
    "span" => %w[class],
    "i" => %w[class],
    "h2" => %w[class],
    "h3" => %w[class]
  }.freeze

  # Markdown rendering of the chapter bodies served by the API.
  def markdown(text)
    options = {
      filter_html: false, # The bodies rely on inline HTML; it is sanitized below.
      link_attributes: { rel: "nofollow", target: "_blank" },
      space_after_headers: true,
      fenced_code_blocks: true
    }

    renderer = Redcarpet::Render::HTML.new(options)
    markdown = Redcarpet::Markdown.new(renderer, autolink: true, footnotes: true)

    html = markdown.render(text.to_s)
    sanitize_content(html).html_safe
  end

  # Strips anything that is not part of the content vocabulary (scripts, event
  # handlers, styles, iframes...).
  def sanitize_content(html)
    Rails::Html::SafeListSanitizer.new.sanitize(
      html,
      tags: MARKDOWN_TAGS,
      attributes: MARKDOWN_ATTRIBUTES
    )
  end


  def render_markdown(text)
    return "".html_safe if text.blank?

    # Decode HTML entities that come escaped from the API
    text = CGI.unescapeHTML(text.to_s)

    # Split the text into lines
    lines = text.split("\n")

    # Initialize a flag to track when to start removing lines
    removing_section = false
    output_lines = []

    lines.each do |line|
      # Check for the first '---' to start removing
      if line.strip == "---"
        removing_section = !removing_section # Toggle the flag
        next # Skip the line with '---'
      end

      # Removing 'Questão x: '
      if line.start_with?("# Questão")
        # Substitute the string with an empty string
        line = line.sub(/Questão \d+:/, "")
      end
      # If we are in the section to remove, skip the lines
      next if removing_section
      # Add the line to output if we are not in the removal section
      output_lines << line
    end
    # Join the output lines back into a single string (preserving blank lines)
    cleaned_text = output_lines.join("\n")

    # Fix for markdown ordered list: replace '1. -' with '1\. -' (escapes dot if matches pattern)
    cleaned_text.gsub!(/^(\d+)\.\s-/, '\1\. -')

    # Escape the dot after a number at the start of a line (unless it's a real list)
    # This prevents Redcarpet from interpreting it as an ordered list
    cleaned_text.gsub!(/^(\s*)(\d+)\.\s(?![\-\d])/, '\1\2\. ')

    # Return the cleaned text
    markdown(cleaned_text)
  end

  # Roman numeral for the part/section counters shown on the index pages.
  def to_roman(number)
    return "" if number.to_i <= 0

    values = [ 1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1 ]
    numerals = [ "M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I" ]

    result = +""
    remaining = number.to_i

    values.each_with_index do |value, index|
      while remaining >= value
        result << numerals[index]
        remaining -= value
      end
    end

    result
  end

  # "A vida em comunidade" -> "A vida em comunidade" (part slugs are internal ids)
  def humanize_slug(slug)
    text = slug.to_s.sub(/\A[\d.]+\-/, "").tr("_", " ")
    text.empty? ? slug.to_s : text[0].upcase + text[1..]
  end

  BOOK_NAMES = {
    "summa_theologiae" => "Suma Teológica",
    "catecismo_pio_x" => "Catecismo de São Pio X",
    "confissoes" => "Confissões",
    "didaque" => "Didaqué"
  }.freeze

  # Part slugs are internal identifiers ("2.1-prima_pars_secundae"), so every known
  # one gets an explicit label instead of being humanized into gibberish.
  PART_NAMES = {
    "summa_theologiae/1-prima_pars" => "Primeira Parte (I)",
    "summa_theologiae/2.1-prima_pars_secundae" => "Segunda Parte, Prima Secundae (I-II)",
    "summa_theologiae/2.2-secundae_secundae" => "Segunda Parte, Secunda Secundae (II-II)",
    "summa_theologiae/3-tertia_pars" => "Terceira Parte (III)",
    "summa_theologiae/4-supplementum" => "Suplemento",
    "summa_theologiae/4.1-supplementum_appendix" => "Apêndice ao Suplemento",
    "catecismo_pio_x/licao_preliminar" => "Lição Preliminar",
    "catecismo_pio_x/primeira_parte" => "Primeira Parte",
    "catecismo_pio_x/segunda_parte" => "Segunda Parte",
    "catecismo_pio_x/terceira_parte" => "Terceira Parte",
    "catecismo_pio_x/quarta_parte" => "Quarta Parte",
    "catecismo_pio_x/quinta_parte" => "Quinta Parte",
    "catecismo_pio_x/apendice" => "Apêndice",
    "didaque/a_instrução_dos_doze_apóstolos" => "Instrução dos Doze Apóstolos",
    "didaque/o_caminho_da_vida_e_o_caminho_da_morte" => "O caminho da vida e o caminho da morte",
    "didaque/a_celebração_litúrgica" => "A celebração litúrgica",
    "didaque/a_vida_em_comunidade" => "A vida em comunidade",
    "didaque/o_fim_dos_tempos" => "O fim dos tempos"
  }.freeze

  CONFISSOES_ORDINALS = {
    "livro_primeiro" => "Primeiro",
    "livro_segundo" => "Segundo",
    "livro_terceiro" => "Terceiro",
    "livro_quarto" => "Quarto",
    "livro_quinto" => "Quinto",
    "livro_sexto" => "Sexto",
    "livro_sétimo" => "Sétimo",
    "livro_oitavo" => "Oitavo",
    "livro_nono" => "Nono",
    "livro_décimo" => "Décimo",
    "livro_décimo-primeiro" => "Décimo Primeiro",
    "livro_décimo-segundo" => "Décimo Segundo",
    "livro_décimo-terceiro" => "Décimo Terceiro"
  }.freeze

  def book_name(book)
    BOOK_NAMES.fetch(book.to_s, book.to_s.humanize)
  end

  # Display label of a part: an explicit name when it is known, "Livro X" for the
  # Confessions and a cleaned-up slug otherwise.
  def part_name(book, part)
    key = "#{book}/#{part}"
    return PART_NAMES[key] if PART_NAMES.key?(key)
    return "Livro #{CONFISSOES_ORDINALS[part]}" if book.to_s == "confissoes" && CONFISSOES_ORDINALS.key?(part)

    humanize_slug(part)
  end
end
