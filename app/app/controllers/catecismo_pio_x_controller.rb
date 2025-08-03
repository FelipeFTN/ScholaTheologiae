class CatecismoPioXController < ApplicationController
  include CatecismoPioXHelper
  include ApplicationHelper

  def index
    @catecismo_parts = catecismo_get_parts
    if @catecismo_parts == nil
      @catecismo_parts = []
    end
    render "books/catecismo_pio_x/index"
  end

  def get_chapters
    @part = params[:part]
    puts "Part: #{@part}"
    @catecismo_chapters = catecismo_get_chapters(@part)
    render "books/catecismo_pio_x/chapters"
  end

  def get_chapter
    @part = params[:part]
    @chapter = params[:chapter]
    data = catecismo_get_chapter(@part, @chapter)
    @content = render_markdown(data)
    render "books/catecismo_pio_x/chapter"
  end
end
