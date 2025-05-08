class Article < ApplicationRecord
  validates :name, presence: true
end
