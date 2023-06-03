package server

import "time"

// CleanupDB removes invalid and expired information from the database. It should be called on server launch.
func CleanupDB() error {
	err := DeleteExpiredSessionsDB()
	if err != nil {
		return err
	}
	return nil
}

// DeleteExpiredSessionsDB deletes all sessions that have expired from the database
func DeleteExpiredSessionsDB() error {
	statement := `DELETE FROM sessions WHERE expires <= $1`
	_, err := db.Exec(statement, time.Now().UTC())
	return err
}
