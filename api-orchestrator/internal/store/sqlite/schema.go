package sqlite

import "database/sql"

func InitSchema(db *sql.DB) error {
	const query = `
CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    integrand TEXT NOT NULL,
    variables_json TEXT NOT NULL,
    evaluations INTEGER NOT NULL,
    status TEXT NOT NULL,
    error_message TEXT,
    result_json TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);`
	_, err := db.Exec(query)
	return err
}
