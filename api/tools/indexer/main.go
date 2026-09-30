// Command indexer builds an FTS5 index of every book's chapter titles and
// bodies into each book's SQLite database.
//
// The chapter bodies live as markdown files (content/<part>/chapter_<n>.md),
// so they are not searchable with plain SQL. This tool reads them and fills a
// per-book FTS5 table named "<book>_fts", tokenized with
// unicode61 remove_diacritics 2 so that searching "graca" finds "graça".
//
// It is idempotent: the FTS table is dropped and rebuilt on every run, so it can
// be re-run safely from the Makefile after the library tarballs are unpacked.
//
// Usage: indexer [-library ./data/library]
package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	htmlTagRe   = regexp.MustCompile(`<[^>]*>`)
	spaceRe     = regexp.MustCompile(`[ \t]+`)
	blankLineRe = regexp.MustCompile(`\n{2,}`)
)

func main() {
	library := flag.String("library", "./data/library", "path to the unpacked book library")
	flag.Parse()

	entries, err := os.ReadDir(*library)
	if err != nil {
		log.Fatalf("indexer: cannot read library at %s: %v", *library, err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		indexed, err := indexBook(filepath.Join(*library, name), name)
		if err != nil {
			log.Fatalf("indexer: %s: %v", name, err)
		}
		log.Printf("indexer: %s: %d chapters indexed", name, indexed)
	}
}

// indexBook rebuilds the FTS table of a single book and returns the number of
// indexed chapters.
func indexBook(dir, name string) (int, error) {
	dbPath := filepath.Join(dir, name+".db")
	if _, err := os.Stat(dbPath); err != nil {
		return 0, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	ftsTable := name + "_fts"
	if _, err := db.Exec("DROP TABLE IF EXISTS " + ftsTable); err != nil {
		return 0, err
	}

	create := `CREATE VIRTUAL TABLE ` + ftsTable + ` USING fts5(
		chapter_title,
		part_title,
		body,
		chapter_number UNINDEXED,
		tokenize = "unicode61 remove_diacritics 2"
	)`
	if _, err := db.Exec(create); err != nil {
		return 0, err
	}

	// The source tables carry one row per article (the chapter title repeats for
	// every article), so chapters are collapsed to one entry per
	// (part, chapter). Rows are read in id order and the last one wins, which is
	// the title of the newest row and matches what the API returns for a part.
	rows, err := db.Query(
		`SELECT part_title, chapter_number, chapter_title, id FROM ` + name +
			` ORDER BY part_title, chapter_number, id`,
	)
	if err != nil {
		return 0, err
	}

	type chapter struct {
		Part   string
		Number int
		Title  string
	}

	collapsed := make(map[string]chapter)
	for rows.Next() {
		var c chapter
		var id int
		if err := rows.Scan(&c.Part, &c.Number, &c.Title, &id); err != nil {
			rows.Close()
			return 0, err
		}
		collapsed[c.Part+"\x00"+strconv.Itoa(c.Number)] = c
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return 0, err
	}
	rows.Close()

	chapters := make([]chapter, 0, len(collapsed))
	for _, c := range collapsed {
		chapters = append(chapters, c)
	}
	sort.Slice(chapters, func(i, j int) bool {
		if chapters[i].Part != chapters[j].Part {
			return chapters[i].Part < chapters[j].Part
		}
		return chapters[i].Number < chapters[j].Number
	})

	insert, err := db.Prepare(`INSERT INTO ` + ftsTable +
		` (chapter_title, part_title, body, chapter_number) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer insert.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	stmt := tx.Stmt(insert)

	missing := 0
	for _, c := range chapters {
		body, err := os.ReadFile(contentPath(dir, c.Part, c.Number))
		if err != nil {
			// A chapter without a body file is still searchable by title.
			missing++
			body = nil
		}
		if _, err := stmt.Exec(c.Title, c.Part, plainText(string(body)), c.Number); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	if missing > 0 {
		log.Printf("indexer: %s: %d chapters have no body file", name, missing)
	}

	return len(chapters), nil
}

// contentPath mirrors the layout written by the scripts/ processors:
// content/<part_title with spaces as underscores>/chapter_<number>.md
func contentPath(dir, part string, number int) string {
	part = strings.ReplaceAll(part, " ", "_")
	return filepath.Join(dir, "content", part, "chapter_"+strconv.Itoa(number)+".md")
}

// plainText removes the markdown decoration that would only add noise to search
// results: HTML tags, heading hashes and the aside the processors keep for
// subtitles.
func plainText(markdown string) string {
	if markdown == "" {
		return ""
	}

	text := htmlTagRe.ReplaceAllString(markdown, " ")
	text = strings.ReplaceAll(text, "\r\n", "\n")

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		trimmed = strings.TrimLeft(trimmed, "#")
		trimmed = strings.TrimPrefix(trimmed, ">")
		lines[i] = trimmed
	}

	text = strings.Join(lines, "\n")
	text = blankLineRe.ReplaceAllString(text, "\n")
	text = spaceRe.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}
