package data

// BOOKS
const (
	SUMMA_THEOLOGIAE = "summa_theologiae"
	PATRISTICA       = "patristica"
	CATECISMO_PIO_X  = "catecismo_pio_x"
)

// SQLiteDBPath is the path to the SQLite database
var SQLiteMAP = map[string]string{
	SUMMA_THEOLOGIAE: "./data/summa_theologiae.db",
	// Not ready yet
	// "patristica":       "./data/patristica.db",
	// "catecismo_pio_x":  "./data/catecismo_pio_x.db",
}
