module ApplicationHelper
  require 'net/http'
  require 'json'

  # Shared API configuration
  BASE_URL = 'http://localhost:8080/v1'.freeze

  def markdown(text)
    options = {
      filter_html: true,
      # hard_wrap: true, # Removed to prevent <br> after each line
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
    # Join the output lines back into a single string (preserving blank lines)
    cleaned_text = output_lines.join("\n")

    # Fix for markdown ordered list: replace '1. -' with '1\\. -' (escapes dot if matches pattern)
    cleaned_text.gsub!(/^(\d+)\.\s-/, '\1\\. -')

    # Escape the dot after a number at the start of a line (unless it's a real list)
    # This prevents Redcarpet from interpreting it as an ordered list
    cleaned_text.gsub!(/^(\s*)(\d+)\.\s(?![\-\d])/, '\1\2\\. ')

    # Return the cleaned text
    markdown(cleaned_text)
  end
end
