module SearchHelper
  require 'net/http'
  require 'json'
  require 'cgi'

  # Where our Go API lives
  API_URL = 'http://localhost:8080/v1'.freeze

  # Main search function - calls the Go API and formats results
  def search_texts(query)
    return [] if query.blank? || query.length < 2

    begin
      # Build the API URL
      encoded_query = CGI.escape(query)
      uri = URI("#{API_URL}/search?q=#{encoded_query}")
      
      # Make the HTTP call with timeouts
      response = Net::HTTP.start(uri.hostname, uri.port, read_timeout: 10, open_timeout: 5) do |http|
        request = Net::HTTP::Get.new(uri)
        request['Content-Type'] = 'application/json'
        http.request(request)
      end

      # Parse results if successful
      if response.code == '200'
        results = JSON.parse(response.body)
        return format_search_results(results)
      else
        Rails.logger.error "Search API failed: #{response.code}"
        return []
      end
      
    rescue => e
      Rails.logger.error "Search error: #{e.message}"
      return []
    end
  end

  private

  # Convert API results to nice format for the view
  def format_search_results(results)
    results.map do |result|
      {
        id: result['id'],
        title: clean_title(result['chapter_title']),
        book: friendly_book_name(result['book']),
        part: friendly_part_name(result['part_title']),
        chapter_number: result['chapter_number'],
        url: build_book_url(result)
      }
    end
  end

  # Clean up messy unicode in titles
  def clean_title(title)
    title.gsub(/\\u([0-9a-f]{4})/i) { [$1.hex].pack('U') }
         .gsub(/\\/, '')
         .strip
  end

  # Convert internal book names to display names
  def friendly_book_name(book)
    case book
    when 'summa_theologiae' then 'Summa Theologica'
    when 'catecismo_pio_x' then 'Catecismo de Pio X'
    else book.humanize
    end
  end

  # Convert part names to Portuguese
  def friendly_part_name(part)
    case part
    when /supplementum/ then 'Suplemento'
    when /primeira_parte/ then 'Primeira Parte'
    when /segunda_parte/ then 'Segunda Parte'
    when /terceira_parte/ then 'Terceira Parte'
    when /quarta_parte/ then 'Quarta Parte'
    else part.humanize
    end
  end

  # Build URLs for each book type
  def build_book_url(result)
    book = result['book']
    part = result['part_title']
    chapter = result['chapter_number']

    case book
    when 'summa_theologiae'
      "/books/summa-theologiae/#{part}/#{chapter}"
    when 'catecismo_pio_x'
      "/books/catecismo-pio-x/#{part}/#{chapter}"
    else
      "/books/"
    end
  end
end
