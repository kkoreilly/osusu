package db

import "time"

// Cleanup removes invalid and expired information from the database. It should be called on server launch.
func Cleanup() error {
	err := DeleteExpiredSessions()
	if err != nil {
		return err
	}
	return nil
}

// DeleteExpiredSessions deletes all sessions that have expired from the database
func DeleteExpiredSessions() error {
	statement := `DELETE FROM sessions WHERE expires <= $1`
	_, err := db.Exec(statement, time.Now().UTC())
	return err
}
