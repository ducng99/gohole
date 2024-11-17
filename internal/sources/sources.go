package sources

type ISource struct {
	ID         int64
	LastFetch  int64
	NumEntries int64
}

const maxStoreSourceSeconds = 3 * 24 * 60 * 60
