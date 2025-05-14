class SummaTheologiaeController < ApplicationController
  include SummaTheologiaeHelper

  def show
    @part = params[:part]
    @question = params[:question]
    @question_data = summa_get_question(@part, @question)
  end

  def search
    @search_term = params[:search_term]
    @search_results = summa_search(@search_term)
  end

  def index
    @summa_parts = summa_get_parts
    if @summa_parts == nil
      @summa_parts = []
    end
    render "books/summa_theologiae/index"
  end

  def get_parts
    @summa_parts = summa_get_parts
    if @summa_parts == nil
      @summa_parts = []
    end
    render "books/summa_theologiae/index"
    # @part = params[:part]
    # @question = params[:question]
    # @question_data = summa_get_question(@part, @question)
  end

  def get_questions
    @part = params[:part]
    @summa_questions = summa_get_questions(@part)
    render "books/summa_theologiae/questions"
  end
end
