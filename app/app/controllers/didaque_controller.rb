class DidaqueController < ApplicationController
  include DidaqueHelper
  include ApplicationHelper

  def index
    @didaque_parts = didaque_get_parts
    if @didaque_parts == nil
      @didaque_parts = []
    end
    
    # Get chapters for each part
    @didaque_books = {}
    @didaque_parts.each do |part|
      @didaque_books[part] = didaque_get_chapters(part)
    end
    
    render "books/didaque/index"
  end

  def get_chapter
    @part = params[:part]
    @chapter = params[:chapter]
    data = didaque_get_chapter(@part, @chapter)
    if data.is_a?(Hash) && data['error']
      redirect_to action: :index
      return
    end
    @content = render_markdown(data)
    render "books/didaque/chapter"
  end
end
