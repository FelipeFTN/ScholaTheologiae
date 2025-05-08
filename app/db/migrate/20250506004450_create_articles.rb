class CreateArticles < ActiveRecord::Migration[8.0]
  def change
    create_table :articles do |t|
      t.string :name
      t.string :author
      t.string :path

      t.timestamps
    end
  end
end
