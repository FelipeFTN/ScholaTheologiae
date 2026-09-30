# frozen_string_literal: true

class SearchController < ApplicationController
  include SearchHelper

  MIN_QUERY_LENGTH = 2

  # Just a placeholder for search page/modal
  def index
  end

  # Handle search requests and show results
  def results
    @query = params[:q]&.strip
    @results = []
    @unavailable = false

    # Only search if we have a decent query
    return if @query.blank? || @query.length < MIN_QUERY_LENGTH

    hits = ApiClient.search(@query)
    if hits.nil?
      @unavailable = true
    else
      @results = format_search_results(hits)
    end
  end
end
