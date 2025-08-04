class BooksController < ApplicationController
  include BooksHelper
  include SummaTheologiaeHelper

  def index
    # This renders the books index view
  end
end
