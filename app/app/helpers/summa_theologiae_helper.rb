module SummaTheologiaeHelper
  include ApplicationHelper

  def summa_get_parts
    begin
      uri = URI("#{BASE_URL}/books/summa_theologiae")
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

  def summa_get_questions(part)
    begin
      uri = URI("#{BASE_URL}/books/summa_theologiae/#{part}")
      response = Net::HTTP.get(uri)

      # Force encoding to UTF-8 to avoid encoding issues
      response = response.force_encoding('UTF-8')
      chapters_hash = JSON.parse(response)
      
      # Convert hash to array of "number: title" strings, similar to Summa format
      chapters_array = chapters_hash.map { |num, title| "#{num}: #{title}" }
      
      # Sort by chapter number
      chapters_array.sort_by { |chapter| chapter.split(':').first.to_i }
    rescue StandardError => e
      puts "Error fetching questions for part #{part}: #{e.message}"
      return []
    end
  end

  def summa_get_question(part, question)
    begin
      uri = URI("#{BASE_URL}/books/summa_theologiae/#{part}/#{question}")
      response = Net::HTTP.get(uri)

      # Force encoding to UTF-8 to avoid encoding issues
      response = response.force_encoding('UTF-8')
      
      # Parse as JSON since the API returns escaped JSON string
      parsed_content = JSON.parse(response)

      return parsed_content
    rescue StandardError => e
      puts "Error fetching question #{question} for part #{part}: #{e.message}"
      return ""
    end
  end
end
