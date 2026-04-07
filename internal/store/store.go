package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

// Poll is a single poll with named options. Options is a JSON-encoded
// array of option strings; Votes is a JSON-encoded map from option name
// to integer vote count. Type is one of: single_choice, multi_choice.
// Status is one of: open, closed, draft.
type Poll struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Options     string `json:"options"` // JSON array as string
	Votes       string `json:"votes"`   // JSON object as string
	Status      string `json:"status"`
	ExpiresAt   string `json:"expires_at"`
	CreatedAt   string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "agora.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS polls(
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT DEFAULT '',
		type TEXT DEFAULT 'single_choice',
		options TEXT DEFAULT '[]',
		votes TEXT DEFAULT '{}',
		status TEXT DEFAULT 'open',
		expires_at TEXT DEFAULT '',
		created_at TEXT DEFAULT(datetime('now'))
	)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_polls_status ON polls(status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_polls_type ON polls(type)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
		resource TEXT NOT NULL,
		record_id TEXT NOT NULL,
		data TEXT NOT NULL DEFAULT '{}',
		PRIMARY KEY(resource, record_id)
	)`)
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

func (d *DB) Create(e *Poll) error {
	e.ID = genID()
	e.CreatedAt = now()
	if e.Status == "" {
		e.Status = "open"
	}
	if e.Type == "" {
		e.Type = "single_choice"
	}
	if e.Options == "" {
		e.Options = "[]"
	}
	if e.Votes == "" {
		e.Votes = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO polls(id, title, description, type, options, votes, status, expires_at, created_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Title, e.Description, e.Type, e.Options, e.Votes, e.Status, e.ExpiresAt, e.CreatedAt,
	)
	return err
}

func (d *DB) Get(id string) *Poll {
	var e Poll
	err := d.db.QueryRow(
		`SELECT id, title, description, type, options, votes, status, expires_at, created_at
		 FROM polls WHERE id=?`,
		id,
	).Scan(&e.ID, &e.Title, &e.Description, &e.Type, &e.Options, &e.Votes, &e.Status, &e.ExpiresAt, &e.CreatedAt)
	if err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []Poll {
	rows, _ := d.db.Query(
		`SELECT id, title, description, type, options, votes, status, expires_at, created_at
		 FROM polls ORDER BY created_at DESC`,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Poll
	for rows.Next() {
		var e Poll
		rows.Scan(&e.ID, &e.Title, &e.Description, &e.Type, &e.Options, &e.Votes, &e.Status, &e.ExpiresAt, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Update(e *Poll) error {
	_, err := d.db.Exec(
		`UPDATE polls SET title=?, description=?, type=?, options=?, votes=?, status=?, expires_at=?
		 WHERE id=?`,
		e.Title, e.Description, e.Type, e.Options, e.Votes, e.Status, e.ExpiresAt, e.ID,
	)
	return err
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM polls WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM polls`).Scan(&n)
	return n
}

// Vote atomically casts a vote for the given option on the given poll.
// Returns the updated votes map. Refuses if the poll is closed, the
// option isn't in the poll's options, or the poll has expired.
//
// This is a read-modify-write transaction inside a single SQL statement
// pair so concurrent voters don't lose increments.
func (d *DB) Vote(pollID, option string) (map[string]int, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var optionsJSON, votesJSON, status, expiresAt string
	err = tx.QueryRow(
		`SELECT options, votes, status, expires_at FROM polls WHERE id=?`,
		pollID,
	).Scan(&optionsJSON, &votesJSON, &status, &expiresAt)
	if err != nil {
		return nil, fmt.Errorf("poll not found")
	}

	if status != "open" {
		return nil, fmt.Errorf("poll is not open")
	}
	if expiresAt != "" && expiresAt < time.Now().UTC().Format(time.RFC3339) {
		return nil, fmt.Errorf("poll has expired")
	}

	var options []string
	if err := json.Unmarshal([]byte(optionsJSON), &options); err != nil {
		return nil, fmt.Errorf("invalid options")
	}
	found := false
	for _, o := range options {
		if o == option {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("option not in poll")
	}

	votes := map[string]int{}
	if votesJSON != "" {
		json.Unmarshal([]byte(votesJSON), &votes)
	}
	votes[option] = votes[option] + 1

	out, _ := json.Marshal(votes)
	_, err = tx.Exec(`UPDATE polls SET votes=? WHERE id=?`, string(out), pollID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return votes, nil
}

// ResetVotes clears all vote counts on a poll while preserving the
// poll's options. Useful for testing or starting a new round.
func (d *DB) ResetVotes(pollID string) error {
	_, err := d.db.Exec(`UPDATE polls SET votes='{}' WHERE id=?`, pollID)
	return err
}

func (d *DB) Search(q string, filters map[string]string) []Poll {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (title LIKE ? OR description LIKE ?)"
		s := "%" + q + "%"
		args = append(args, s, s)
	}
	if v, ok := filters["type"]; ok && v != "" {
		where += " AND type=?"
		args = append(args, v)
	}
	if v, ok := filters["status"]; ok && v != "" {
		where += " AND status=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(
		`SELECT id, title, description, type, options, votes, status, expires_at, created_at
		 FROM polls WHERE `+where+`
		 ORDER BY created_at DESC`,
		args...,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Poll
	for rows.Next() {
		var e Poll
		rows.Scan(&e.ID, &e.Title, &e.Description, &e.Type, &e.Options, &e.Votes, &e.Status, &e.ExpiresAt, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

// Stats returns aggregate metrics: total polls, total votes cast across
// all polls, by_status and by_type breakdowns.
func (d *DB) Stats() map[string]any {
	m := map[string]any{
		"total":       d.Count(),
		"total_votes": 0,
		"by_status":   map[string]int{},
		"by_type":     map[string]int{},
	}

	// Sum vote counts by parsing each poll's votes JSON
	rows, _ := d.db.Query(`SELECT votes FROM polls`)
	if rows != nil {
		defer rows.Close()
		total := 0
		for rows.Next() {
			var vj string
			rows.Scan(&vj)
			var v map[string]int
			if json.Unmarshal([]byte(vj), &v) == nil {
				for _, c := range v {
					total += c
				}
			}
		}
		m["total_votes"] = total
	}

	if rows, _ := d.db.Query(`SELECT status, COUNT(*) FROM polls GROUP BY status`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_status"] = by
	}

	if rows, _ := d.db.Query(`SELECT type, COUNT(*) FROM polls GROUP BY type`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_type"] = by
	}

	return m
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
