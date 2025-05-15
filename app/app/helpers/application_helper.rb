module ApplicationHelper
  def markdown(text)
    options = {
      filter_html: true,
      hard_wrap: true,
      link_attributes: { rel: 'nofollow', target: "_blank" },
      space_after_headers: true,
      fenced_code_blocks: true
    }

    renderer = Redcarpet::Render::HTML.new(options)
    markdown = Redcarpet::Markdown.new(renderer, extensions = {})

    markdown.render(text).html_safe
  end

  def render_markdown(text)
    # Split the text into lines
    lines = text.split("\n")

    # Initialize a flag to track when to start removing lines
    removing_section = false
    output_lines = []

    lines.each do |line|
      # Check for the first '---' to start removing
      if line.strip == '---'
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
    # Join the output lines back into a single string
    cleaned_text = output_lines.join("\n")

    # Return the cleaned text
    markdown(cleaned_text)
  end
end
