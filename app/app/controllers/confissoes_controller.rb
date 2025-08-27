class ConfissoesController < ApplicationController
  include ConfissoesHelper
  include ApplicationHelper

  def index
    @confissoes_parts = confissoes_get_parts
    if @confissoes_parts == nil
      @confissoes_parts = []
    end
    
    # Get chapters for each part
    @confissoes_books = {}
    @confissoes_parts.each do |part|
      @confissoes_books[part] = confissoes_get_chapters(part)
    end
    
    render "books/confissoes/index"
  end

  def get_chapter
    @part = params[:part]
    @chapter = params[:chapter]
    data = confissoes_get_chapter(@part, @chapter)
    if data.is_a?(Hash) && data['error']
      redirect_to action: :index
      return
    end
    @content = render_markdown(data)
    render "books/confissoes/chapter"
  end
end
