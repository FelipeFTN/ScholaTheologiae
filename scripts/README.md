# 📁 Standard for books databases

## SQLite table schema

That's quite important, every single book should have this structure for the API to handle and search it correctly.
Even the database structure... one addition to a singular book that doesn't contemplate the others, should generate quite a mess.

```sql
CREATE TABLE IF NOT EXISTS book_name (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    part_title TEXT NOT NULL,
    chapter_title TEXT NOT NULL,
    chapter_number INTEGER NOT NULL
    article_title TEXT,
    article_number INTEGER
);
```

*All the number fields are necessary for URL search, for example, in front-end, it will use paths like this `/books/book_name/part_title/chapter_number/article_number`.*
***Why not using part_number too? R. Because not; seems ok like this.***

#### The book content should be stored in a separate file with the following structure standard:

```
book_name/
    part_title/
        chapter_n.md
        chapter_2.md
        chapter_3.md
        chapter_4.md
        chapter_5.md
        [...]
```

It's also important to note, this folder `book_name` should be **tar.gziped** and the output files should be placed at [books](../books/).
During the API's building & compilation process, the Makefile will build all the books inside of it, and insert it into [library](../api/data/library/)! :pray:
