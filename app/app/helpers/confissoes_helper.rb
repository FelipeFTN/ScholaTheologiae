module ConfissoesHelper
  include ApplicationHelper

  def confissoes_get_parts
    begin
      uri = URI("#{BASE_URL}/books/confissoes")
      response = Net::HTTP.get(uri)
      if response.empty?
        puts "No response from server"
        return []
      end
      # Force encoding to UTF-8 to avoid encoding issues
      response = response.force_encoding('UTF-8')
      return JSON.parse(response)
    rescue StandardError => e
      puts "Error fetching parts: #{e.message}"
      return []
    end
  end

  def confissoes_get_chapters(part)
    begin
      # Encode the part to ensure it is ASCII only
      encoded_part = URI.encode_www_form_component(part)
      uri = URI("#{BASE_URL}/books/confissoes/#{encoded_part}")
      response = Net::HTTP.get(uri)
      # Force encoding to UTF-8 to avoid encoding issues
      response = response.force_encoding('UTF-8')
      chapters_hash = JSON.parse(response)

      # Convert hash to array of "number: title" strings, similar to format
      chapters_array = chapters_hash.map { |num, title| "#{num}: #{title}" }

      # Sort by chapter number
      chapters_array.sort_by { |chapter| chapter.split(':').first.to_i }
    rescue StandardError => e
      puts "Error fetching chapters for part #{part}: #{e.message}"
      return []
    end
  end

  def confissoes_get_chapter(part, chapter)
    begin
      # Encode the part to ensure it is ASCII only
      encoded_part = URI.encode_www_form_component(part)
      uri = URI("#{BASE_URL}/books/confissoes/#{encoded_part}/#{chapter}")
      response = Net::HTTP.get(uri)
      # Force encoding to UTF-8 to avoid encoding issues
      response = response.force_encoding('UTF-8')
      
      # Parse as JSON since the API returns escaped JSON string
      parsed_content = JSON.parse(response)
      return parsed_content
    rescue StandardError => e
      puts "Error fetching chapter #{chapter} for part #{part}: #{e.message}"
      return ""
    end
  end
end
