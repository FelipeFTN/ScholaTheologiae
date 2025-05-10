Rails.application.routes.draw do
  get "health" => "rails/health#show", as: :health_check

  get "home" => "home#index", as: :home

  get "/books" => "books#index", as: :books_all

  get "/books/:name" => "books#get", as: :books_get

  # Summa Theologiae specific routes because of the way the data is structured
  get "/books/summa-theologiae" => "books#summa_theologiae", as: :books_summa_theologiae

  root "home#index"
end
