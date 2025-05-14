Rails.application.routes.draw do
  get "health" => "rails/health#show", as: :health_check

  get "home" => "home#index", as: :home

  get "books" => "books#index", as: :books

  # Summa Theologiae specific routes because of the way the data is structured
  get "/books/summa-theologiae" => "summa_theologiae#get_parts"
  get "/books/summa-theologiae/:part" => "summa_theologiae#get_questions"
  # get "/books/summa-theologiae/:part/:question" => "summa_theologiae_question", as: :summa_theologiae_question

  # Catch all other routes and return 404
  match '*path', to: ->(env) { [ 404, { 'Content-Type' => 'text/html' }, [ File.read(Rails.root.join('public', '404.html')) ] ] }, via: :all

  root "home#index"
end
