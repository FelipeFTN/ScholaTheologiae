class ArticlesController < ApplicationController
  def index
    @articles = Article.all
  end

  def get
    @article = Article.find_by(path: params[:path])
  end
end
