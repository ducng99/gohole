package db

import (
	"database/sql"
	"errors"

	"github.com/ducng99/gohole/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

type HoleDB struct {
	Conn *sql.DB
}

type dbConnection interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}

type ISource struct {
	ID         int64
	Url        string
	LastFetch  int64
	NumEntries int64
}

type IDomainEntry struct {
	Domain   string
	SourceID int64
}

var db *HoleDB = nil

var ErrSourceNotFound = errors.New("Source not found")

func New() *HoleDB {
	if db == nil {
		var err error

		_db, err := sql.Open("sqlite3", "file:gohole.db")
		if err != nil {
			logger.Fatalf("%v", err)
		}

		db = &HoleDB{
			Conn: _db,
		}
	}

	db.initDB()

	return db
}

func (db *HoleDB) Close() {
	db.Conn.Close()
}

func (db *HoleDB) initDB() {
	_, err := db.Conn.Exec(`CREATE TABLE IF NOT EXISTS hole_sources (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		last_fetch INTEGER DEFAULT (0) NOT NULL,
		num_entries INTEGER DEFAULT (0) NOT NULL
	)`)
	if err != nil {
		logger.Fatalf("Could not create `hole_sources` table.\n%v\n", err)
	}

	_, err = db.Conn.Exec(`CREATE TABLE IF NOT EXISTS hole_entries (
		entry_domain TEXT NOT NULL,
		source_id INTEGER NOT NULL
	)`)
	if err != nil {
		logger.Fatalf("Could not create `hole_entries` table.\n%v\n", err)
	}
}

func (db *HoleDB) Vacuum() error {
	_, err := db.Conn.Exec("VACUUM")
	if err != nil {
		return err
	}

	return nil
}

func AddSource(db dbConnection, url string) (int64, error) {
	result, err := db.Exec(`INSERT INTO hole_sources (url) VALUES (?)`, url)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateSource(db dbConnection, id int64, numEntries int) error {
	_, err := db.Exec(
		`UPDATE hole_sources
		SET last_fetch = unixepoch(),
			num_entries = ?
		WHERE id = ?`,
		numEntries, id,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetSourceID(db dbConnection, url string) (int64, error) {
	row := db.QueryRow("SELECT id FROM hole_sources WHERE url = ?", url)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var id int64 = 0
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return id, nil
}

func GetSourceUrl(db dbConnection, sourceID int64) (string, error) {
	row := db.QueryRow("SELECT url FROM hole_sources WHERE id = ?", sourceID)
	if err := row.Err(); err != nil {
		return "", err
	}

	var url string
	if err := row.Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrSourceNotFound
		}

		return "", err
	}

	return url, nil
}

func GetSourceStats(db dbConnection, sourceID int64) (lastFetch int64, numEntries int64, err error) {
	row := db.QueryRow("SELECT last_fetch, num_entries FROM hole_sources WHERE id = ?", sourceID)
	if err := row.Err(); err != nil {
		return 0, 0, err
	}

	if err := row.Scan(&lastFetch, &numEntries); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, ErrSourceNotFound
		}

		return 0, 0, err
	}

	return lastFetch, numEntries, nil
}

func GetSources(db dbConnection) ([]ISource, error) {
	rows, err := db.Query("SELECT id, url, last_fetch, num_entries FROM hole_sources")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sources := make([]ISource, 0, 5)

	for rows.Next() {
		var id int64
		var url string
		var lastFetch int64
		var numEntries int64

		if err = rows.Scan(&id, &url, &lastFetch, &numEntries); err != nil {
			return nil, err
		}

		sources = append(sources, ISource{ID: id, Url: url, LastFetch: lastFetch, NumEntries: numEntries})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}

func RemoveSource(db dbConnection, sourceID int64) error {
	_, err := db.Exec("DELETE FROM hole_sources WHERE id = ?", sourceID)
	if err != nil {
		return err
	}

	return nil
}

func AddDomains(db dbConnection, domains []string, sourceID int64) error {
	addStmt, err := db.Prepare("INSERT INTO hole_entries (entry_domain, source_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer addStmt.Close()

	for _, domain := range domains {
		_, err = addStmt.Exec(domain, sourceID)
		if err != nil {
			return err
		}
	}

	return nil
}

func ClearSourceDomains(db dbConnection, sourceID int64) error {
	_, err := db.Exec("DELETE FROM hole_entries WHERE source_id = ?", sourceID)
	if err != nil {
		return err
	}

	return nil
}

func GetDomains(db dbConnection) ([]IDomainEntry, error) {
	rows, err := db.Query("SELECT entry_domain, source_id FROM hole_entries GROUP BY entry_domain")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	domains := make([]IDomainEntry, 0, 300000)

	for rows.Next() {
		domain := IDomainEntry{}

		if err = rows.Scan(&domain.Domain, &domain.SourceID); err != nil {
			return nil, err
		}

		domains = append(domains, domain)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return domains, nil
}
