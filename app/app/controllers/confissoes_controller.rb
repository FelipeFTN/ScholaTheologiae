class ConfissoesController < ApplicationController
  include BookReader

  def book
    "confissoes"
  end

  def book_path
    "/books/confissoes"
  end

  def book_views
    "books/confissoes"
  end
end
