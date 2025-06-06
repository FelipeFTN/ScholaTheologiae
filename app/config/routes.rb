Rails.application.routes.draw do
  get "health" => "rails/health#show", as: :health_check

  get "home" => "home#index", as: :home

  get "books" => "books#index", as: :books

  get "divine_office" => "divine_office#index", as: :divine_office
  get "calendar" => "calendar#index", as: :calendar
  get "daily_liturgy" => "daily_liturgy#index", as: :daily_liturgy
  get "daily_meditation" => "daily_meditation#index", as: :daily_meditation
  get "caminho" => "caminho#index", as: :caminho
  get "via_sacra" => "via_sacra#index", as: :via_sacra

  # Summa Theologiae specific routes because of the way the data is structured
  get "/books/summa-theologiae" => "summa_theologiae#get_parts"

  # This :part might contain a dot, so we need to use a regex to match it
  get "/books/summa-theologiae/:part" => "summa_theologiae#get_questions", constraints: { part: /[^\/]+/ }

  get "/books/summa-theologiae/:part/:question" => "summa_theologiae#get_question", constraints: { part: /[^\/]+/, question: /[^\/]+/ }

  # Catch all other routes and return 404
  match '*path', to: ->(env) { [ 404, { 'Content-Type' => 'text/html' }, [ File.read(Rails.root.join('public', '404.html')) ] ] }, via: :all

  # root "home#index" # We are not using this page yet
  root "books#index"
end
