class SummaTheologiaeController < ApplicationController
  include SummaTheologiaeHelper
  include ApplicationHelper

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
    puts "Part: #{@part}"
    @summa_questions = summa_get_questions(@part)
    render "books/summa_theologiae/questions"
  end

  def get_question
    @part = params[:part]
    @question = params[:question]
    data = summa_get_question(@part, @question)
    # Split the text into lines
    lines = data.split("\n")

    # Initialize a flag to track when to start removing lines
    removing_section = false
    output_lines = []

    lines.each_with_index do |line, index|
      # Check for the first '---' to start removing
      if line.strip == '---'
        removing_section = !removing_section # Toggle the flag
        next # Skip the line with '---'
      end

      # If we are in the section to remove, skip the lines
      next if removing_section

      # Add the line to output if we are not in the removal section
      output_lines << line
    end

    # Join the output lines back into a single string
    cleaned_text = output_lines.join("\n")

    # @question_content = @question_content.gsub("\n", "<br>")
    @content = markdown(cleaned_text)

    render "books/summa_theologiae/question"
  end
end
