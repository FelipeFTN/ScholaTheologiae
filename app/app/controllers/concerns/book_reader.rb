# frozen_string_literal: true

# Shared behaviour of the books whose content comes straight from the API: the
# index lists every part with its chapters, and a chapter page renders the
# markdown body with links to the neighbouring chapters.
#
# Including controllers must define #book (the API/database name) and #book_path
# (the public route prefix), plus #book_views for the template folder.
module BookReader
  extend ActiveSupport::Concern

  included do
    include ApplicationHelper
  end

  def index
    @parts = ApiClient.parts(book)
    @chapters = chapters_by_part(@parts)

    render "#{book_views}/index"
  end

  def get_chapter
    @part = params[:part]
    @chapter = params[:chapter].to_s
    @chapters = ApiClient.chapters(book, @part)
    @content = ApiClient.chapter(book, @part, @chapter)

    if @content.nil?
      flash[:alert] = t_unavailable_chapter
      redirect_to book_path
      return
    end

    @content = render_markdown(@content)
    set_chapter_navigation

    render "#{book_views}/chapter"
  end

  private

  # Chapters of every part, keyed by part: { "part_title" => [{ number:, title: }] }
  def chapters_by_part(parts)
    parts.to_a.each_with_object({}) do |part, chapters|
      chapters[part] = ApiClient.chapters(book, part)
    end
  end

  # Chapter links are derived from the real chapter list, so the last chapter no
  # longer offers a "next" that does not exist and the header can show progress.
  # Progress is only shown when the chapter numbering is not 1..N for the part
  # (Didaqué numbers its chapters across the whole work, for instance).
  def set_chapter_navigation
    numbers = @chapters.to_a.map { |chapter| chapter[:number] }
    position = numbers.index(@chapter.to_i)

    @chapter_total = numbers.size
    @chapter_progress = if position && position + 1 != @chapter.to_i
      "#{position + 1} de #{numbers.size}"
    end
    @previous_chapter = position && position.positive? ? numbers[position - 1] : nil
    @next_chapter = position ? numbers[position + 1] : nil
  end

  def t_unavailable_chapter
    "Não encontramos este capítulo. Voltamos para a lista de capítulos."
  end
end
