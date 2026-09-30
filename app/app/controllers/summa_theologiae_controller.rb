class SummaTheologiaeController < ApplicationController
  include ApplicationHelper

  BOOK = "summa_theologiae"
  BOOK_PATH = "/books/summa-theologiae"

  def index
    get_parts
  end

  # Parts are numbered ("1-prima_pars", "2.1-prima_pars_secundae"), so they are
  # sorted numerically instead of alphabetically. nil means the API did not answer.
  def get_parts
    parts = ApiClient.parts(BOOK)
    @parts = parts&.sort_by { |part| part.split(".").map(&:to_i) }

    render "books/summa_theologiae/index"
  end

  def get_questions
    @part = params[:part]
    @questions = ApiClient.chapters(BOOK, @part)

    if @questions.nil?
      @not_found = true
      render "books/summa_theologiae/questions", status: :not_found
      return
    end

    render "books/summa_theologiae/questions"
  end

  def get_question
    @part = params[:part]
    @question = params[:question].to_s
    @questions = ApiClient.chapters(BOOK, @part)
    @content = ApiClient.chapter(BOOK, @part, @question)

    if @content.nil?
      flash[:alert] = "Não encontramos esta questão. Voltamos para a lista de questões."
      redirect_to "#{BOOK_PATH}/#{@part}"
      return
    end

    @content = render_markdown(@content)
    set_question_navigation

    render "books/summa_theologiae/question"
  end

  private

  # Question links come from the real question list, so the last question no
  # longer offers a "next" that does not exist. Progress is only shown when the
  # question numbering is not 1..N for the part.
  def set_question_navigation
    numbers = @questions.to_a.map { |question| question[:number] }
    position = numbers.index(@question.to_i)

    @question_total = numbers.size
    @question_progress = if position && position + 1 != @question.to_i
      "#{position + 1} de #{numbers.size}"
    end
    @previous_question = position && position.positive? ? numbers[position - 1] : nil
    @next_question = position ? numbers[position + 1] : nil
  end
end
