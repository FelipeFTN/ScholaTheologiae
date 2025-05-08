Rails.application.routes.draw do
  get "health" => "rails/health#show", as: :health_check

  get "home" => "home#index", as: :home

  get "/articles" => "articles#index", as: :articles_all

  get "/articles/:path" => "articles#get", as: :articles_get

  root "home#index"
end
