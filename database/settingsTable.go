package database

import (
	"database/sql"
	"log"
)

type SettingsTable struct {
	pkDb *sql.DB
}

func (t *SettingsTable) UpdateDefault(guildID string, channelID string) error {
	tx, err := t.pkDb.Begin()
	if err != nil {
		return err
	}
	stmt, err := t.pkDb.Prepare(`INSERT OR REPLACE INTO settings (guildID, defaultChannel) VALUES (?, ?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(guildID, channelID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (t *SettingsTable) QueryDefault(guildID string) (string, error) {
	stmt, err := t.pkDb.Prepare(`SELECT defaultChannel
								 FROM settings 
								 WHERE guildID = ?`)
	if err != nil {
		return "", err
	}
	result := ""
	if err := stmt.QueryRow(guildID).Scan(&result); err != nil {
		if err != sql.ErrNoRows {
			log.Println("Error while searching ", err)
			return "", err
		}
	}
	return result, nil
}
