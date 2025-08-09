package data

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
