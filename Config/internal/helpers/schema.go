package helpers

import "database/sql"

func InitSchema(db *sql.DB) error {
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS agendas (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			uca_id VARCHAR(255)
		);`,
		`CREATE TABLE IF NOT EXISTS alerts (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			agenda_id VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			FOREIGN KEY (agenda_id) REFERENCES agendas(id)
			ON DELETE CASCADE
		);`,
	}

	for _, schema := range schemas {
		if _, err := db.Exec(schema); err != nil {
			return err
		}
	}

	return nil
}
