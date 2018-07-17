package database

import (
	"database/sql"
	"log"
)

type AliasTable struct {
	pkDb *sql.DB
}

type Alias struct {
	GuildID   string
	ChannelID string
	Alias     string
}

func (t *AliasTable) InsertUpdateAlias(guildID string, channelID string, alias string) error {
	tx, err := t.pkDb.Begin()
	if err != nil {
		log.Println("Error while inserting aliases: ", err)
		return err
	}
	stmt, err := tx.Prepare(`SELECT COUNT(*) FROM alias WHERE guildID = ? AND channelID = ?`)
	if err != nil {
		log.Println("Error while inserting aliases: ", err)
		tx.Rollback()
		return err
	}
	count := 0
	err = stmt.QueryRow(guildID, channelID).Scan(&count)
	if err != nil {
		log.Println("Error while inserting aliases: ", err)
		tx.Rollback()
		return err
	}
	if count == 0 {
		stmt, err = tx.Prepare(`INSERT INTO alias (guildID, channelID, alias) VALUES (?,?,?)`)
		if err != nil {
			log.Println("Error while inserting aliases: ", err)
			tx.Rollback()
			return err
		}
		_, err = stmt.Exec(guildID, channelID, alias)
		if err != nil {
			log.Println("Error while inserting aliases: ", err)
			tx.Rollback()
			return err
		}
	} else {
		stmt, err = tx.Prepare(`UPDATE alias SET alias = ?`)
		if err != nil && err != sql.ErrNoRows {
			log.Println("Error while inserting aliases: ", err)
			tx.Rollback()
			return err
		}
		_, err = stmt.Exec(alias)
		if err != nil {
			log.Println("Error while inserting aliases: ", err)
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (t *AliasTable) QueryAliases() ([]Alias, error) {
	stmt, err := t.pkDb.Prepare(`SELECT guildID, channelID, alias
								 FROM alias`)
	if err != nil {
		log.Println("Error while retrieving aliases: ", err)
		return nil, err
	}
	result := []Alias{}
	rows, err := stmt.Query()
	for rows.Next() {
		alias := Alias{}
		err = rows.Scan(&alias.GuildID, &alias.ChannelID, &alias.Alias)
		if err != nil {
			log.Println("Error while retrieving aliases: ", err)
			return nil, err
		}
		result = append(result, alias)
	}
	return result, nil
}
