module BooksHelper
  require 'net/http'
  require 'json'

  BASE_URL = 'http://localhost:8080/v1'.freeze

  def summa_get_parts
    begin
      uri = URI("#{BASE_URL}/summa-theologiae")
      response = Net::HTTP.get(uri)
      if response.empty?
        puts "No response from server"
        return []
      end
      return JSON.parse(response)
    rescue StandardError => e
      puts "Error fetching parts: #{e.message}"
      return []
    end
  end

  def summa_get_questions(part)
    begin
      uri = URI("#{BASE_URL}/summa-theologiae#{part}")
      response = Net::HTTP.get(uri)
      JSON.parse(response)
    rescue StandardError => e
      puts "Error fetching questions for part #{part}: #{e.message}"
      return []
    end
  end 

  def summa_get_question(part, question)
    begin
      uri = URI("#{BASE_URL}/summa-theologiae#{part}/#{question}")
      response = Net::HTTP.get(uri)
      JSON.parse(response)
    rescue StandardError => e
      puts "Error fetching question #{question} for part #{part}: #{e.message}"
      return ""
    end
  end
end
