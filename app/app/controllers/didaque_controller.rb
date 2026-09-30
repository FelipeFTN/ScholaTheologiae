class DidaqueController < ApplicationController
  include BookReader

  def book
    "didaque"
  end

  def book_path
    "/books/didaque"
  end

  def book_views
    "books/didaque"
  end
end
