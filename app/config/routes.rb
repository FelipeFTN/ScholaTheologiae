Rails.application.routes.draw do
  get "health" => "rails/health#show", as: :health_check

  get "home" => "home#index", as: :home

  get "books" => "books#index", as: :books

  # Summa Theologiae specific routes because of the way the data is structured
  get "/books/summa-theologiae" => "books#summa_theologiae", as: :books_summa_theologiae

  root "home#index"

  # Catch all other routes and return 404
  match '*path', to: ->(env) { [404, {'Content-Type' => 'text/html'}, [File.read(Rails.root.join('public', '404.html'))]] }, via: :all
end
