package utils

import "database/sql"

func GetDBConnection(DSN string) (*sql.DB, error) { //conf DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	return db, nil
}
