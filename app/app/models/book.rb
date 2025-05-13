class Book < ApplicationRecord
  validates :name, presence: true
end
