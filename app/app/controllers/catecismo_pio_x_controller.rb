class CatecismoPioXController < ApplicationController
  include BookReader

  def book
    "catecismo_pio_x"
  end

  def book_path
    "/books/catecismo-pio-x"
  end

  def book_views
    "books/catecismo_pio_x"
  end
end
