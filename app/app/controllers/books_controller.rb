class BooksController < ApplicationController
  include BooksHelper

  def index
    @summa_parts = summa_get_parts
    if @summa_parts == nil
      @summa_parts = []
    end
  end

  def summa_theologiae
    @summa_parts = summa_get_parts
    if @summa_parts == nil
      @summa_parts = []
    end
    @part = params[:part]
    @question = params[:question]
    @question_data = summa_get_question(@part, @question)
  end

  def get
    @name = params[:name]
    # Add any additional logic needed to fetch and display the specific book content
  end
end
