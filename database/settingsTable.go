package database

import (
	"database/sql"
	"log"
)

type SettingsTable struct {
	pkDb *sql.DB
}

type Settings struct {
	DefaultChannel string
	ArchiveChannel string
}

func (t *SettingsTable) UpdateDefault(guildID string, defaultChannelID string, archiveChannelID string) error {
	tx, err := t.pkDb.Begin()
	if err != nil {
		return err
	}
	stmt, err := t.pkDb.Prepare(`INSERT OR REPLACE INTO settings (guildID, defaultChannel, archiveChannel) VALUES (?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(guildID, defaultChannelID, archiveChannelID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (t *SettingsTable) QueryDefault(guildID string) (Settings, error) {
	stmt, err := t.pkDb.Prepare(`SELECT defaultChannel, archiveChannel
								 FROM settings 
								 WHERE guildID = ?`)
	if err != nil {
		return Settings{}, err
	}
	result := Settings{}
	if err := stmt.QueryRow(guildID).Scan(&result.DefaultChannel, &result.ArchiveChannel); err != nil {
		if err != sql.ErrNoRows {
			log.Println("Error while searching ", err)
			return Settings{}, err
		}
	}
	return result, nil
}
