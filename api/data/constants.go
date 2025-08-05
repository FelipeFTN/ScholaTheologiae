package data

// BOOKS
const (
	SUMMA_THEOLOGIAE = "tmp_summa"
	PATRISTICA       = "patristica"
	CATECISMO_PIO_X  = "catecismo_pio_x"
)

// SQLiteDBPath is the path to the SQLite database
var DATABASES = map[string]string{
	SUMMA_THEOLOGIAE: "./data/library/summa_theologiae/summa_theologiae.db",
	CATECISMO_PIO_X:  "./data/library/catecismo_pio_x/catecismo_pio_x.db",
	// Not ready yet
	// "patristica":       "./data/patristica.db",
}
