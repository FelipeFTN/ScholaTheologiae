class CatecismoPioXController < ApplicationController
  include CatecismoPioXHelper
  include ApplicationHelper

  def index
    @catecismo_parts = catecismo_get_parts
    if @catecismo_parts == nil
      @catecismo_parts = []
    end
    
    # Get chapters for each part
    @catecismo_books = {}
    @catecismo_parts.each do |part|
      @catecismo_books[part] = catecismo_get_chapters(part)
    end
    
    render "books/catecismo_pio_x/index"
  end

  def get_chapter
    @part = params[:part]
    @chapter = params[:chapter]
    data = catecismo_get_chapter(@part, @chapter)
    if data.is_a?(Hash) && data['error']
      redirect_to action: :index, part: @part
      return
    end
    @content = render_markdown(data)
    render "books/catecismo_pio_x/chapter"
  end
end
