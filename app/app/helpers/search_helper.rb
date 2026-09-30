# frozen_string_literal: true

# Turns the raw hits returned by ApiClient.search into the shape the results page
# needs: a display location and a link.
module SearchHelper
  include ApplicationHelper

  # Public route of every book the API knows about.
  BOOK_URLS = {
    "summa_theologiae" => "/books/summa-theologiae",
    "catecismo_pio_x" => "/books/catecismo-pio-x",
    "confissoes" => "/books/confissoes",
    "didaque" => "/books/didaque"
  }.freeze

  # Convert API hits into the entries the results page renders. Hits without a
  # public route are dropped instead of linking somewhere unrelated.
  def format_search_results(hits)
    hits.filter_map do |hit|
      url = book_url(hit)
      next if url.nil?

      {
        title: hit[:title].to_s.strip,
        snippet: hit[:snippet],
        book: book_name(hit[:book]),
        part: part_name(hit[:book], hit[:part]),
        chapter_number: hit[:chapter],
        url: url
      }
    end
  end

  # nil when the book has no public route or the hit lacks a location, so the
  # caller can drop it instead of building a broken link.
  def book_url(hit)
    base = BOOK_URLS[hit[:book].to_s]
    return nil if base.nil? || hit[:part].blank? || hit[:chapter].blank?

    "#{base}/#{hit[:part]}/#{hit[:chapter]}"
  end
end
