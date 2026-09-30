# frozen_string_literal: true

require "net/http"
require "json"
require "erb"

# Thin client for the Go API that backs every book, chapter and search page.
#
# Every request is bounded by timeouts, every path segment is URL-encoded and
# failures are logged and reported as nil, so the callers can render a proper
# "conteúdo indisponível" state instead of an empty page.
module ApiClient
  BASE_URL = ENV.fetch("SCHOLA_API_URL", "http://localhost:8080/v1").freeze
  OPEN_TIMEOUT = 5
  READ_TIMEOUT = 10

  class << self
    # Parts (sections, books) of one book. Returns an array of strings or nil.
    def parts(book)
      payload = get("/books/#{segment(book)}")
      payload.is_a?(Array) ? payload : nil
    end

    # Chapters of one part, sorted by number: [{ number:, title: }, ...] or nil.
    def chapters(book, part)
      payload = get("/books/#{segment(book)}/#{segment(part)}")
      return nil unless payload.is_a?(Hash)

      payload
        .map { |number, title| { number: number.to_i, title: title.to_s } }
        .sort_by { |chapter| chapter[:number] }
    end

    # Markdown body of one chapter, or nil.
    def chapter(book, part, number)
      payload = get("/books/#{segment(book)}/#{segment(part)}/#{segment(number)}")
      payload.is_a?(String) && payload.present? ? payload : nil
    end

    # Search hits: [{ book:, part:, chapter:, title:, snippet: }, ...] or nil.
    def search(query)
      payload = get("/search", q: query)
      return nil unless payload.is_a?(Array)

      payload.map { |result| hit(result) }
    end

    private

    def get(path, params = {})
      uri = URI("#{BASE_URL}#{path}")
      uri.query = URI.encode_www_form(params) if params.any?

      response = perform(uri)

      case response.code.to_i
      when 200
        JSON.parse(response.body.force_encoding("UTF-8"))
      when 404
        Rails.logger.info("API: #{uri.path} has no content")
        nil
      else
        Rails.logger.error("API: #{uri.path} answered #{response.code}")
        nil
      end
    rescue StandardError => e
      Rails.logger.error("API: #{uri.path} failed: #{e.class}: #{e.message}")
      nil
    end

    def perform(uri)
      Net::HTTP.start(uri.hostname, uri.port, open_timeout: OPEN_TIMEOUT, read_timeout: READ_TIMEOUT) do |http|
        request = Net::HTTP::Get.new(uri.request_uri)
        request["Content-Type"] = "application/json"
        http.request(request)
      end
    end

    def hit(result)
      {
        book: result["book"],
        part: result["part_title"],
        chapter: result["chapter_number"],
        title: result["chapter_title"].to_s,
        snippet: result["snippet"]
      }
    end

    # Percent-encode a path segment (spaces become %20, not "+").
    def segment(value)
      ERB::Util.url_encode(value.to_s)
    end
  end
end
