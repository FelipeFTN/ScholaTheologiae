package data

import "strings"

type Data struct {
	SQLite  *SQLite
	Library *Library
}

func New() (*Data, error) {
	sqlite, err := NewSQLite(DATABASES)
	if err != nil {
		return nil, err
	}

	return &Data{
		SQLite: sqlite,
	}, nil
}

func removeAccents(input string) string {
	replacements := map[string]string{
		"á": "a", "à": "a", "ã": "a", "â": "a", "ä": "a",
		"é": "e", "è": "e", "ê": "e", "ë": "e",
		"í": "i", "ì": "i", "î": "i", "ï": "i",
		"ó": "o", "ò": "o", "õ": "o", "ô": "o", "ö": "o",
		"ú": "u", "ù": "u", "û": "u", "ü": "u",
		"ç": "c", "ñ": "n",
	}

	for accented, unaccented := range replacements {
		input = strings.ReplaceAll(input, accented, unaccented)
	}
	return input
}
