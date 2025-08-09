# frozen_string_literal: true

class SearchController < ApplicationController
  include SearchHelper

  # Just a placeholder for search page/modal
  def index
  end

  # Handle search requests and show results
  def results
    @query = params[:q]&.strip
    @results = []
    
    # Only search if we have a decent query
    if @query.present? && @query.length >= 2
      @results = search_texts(@query)
      Rails.logger.info "Search: '#{@query}' found #{@results.length} results"
    end
  end
end
