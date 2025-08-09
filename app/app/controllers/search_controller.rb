# frozen_string_literal: true

class SearchController < ApplicationController
  def index
    # This will be the search page/modal
  end

  def results
    # This will handle the search results
    @query = params[:q]
    # For now, we'll just render the results without calling the API
    # We'll implement the API call later
  end
end
