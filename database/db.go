package database

import (
	"database/sql"
	"time"
)

type Database struct {
	pkDb     *sql.DB
	Image    *ImageTable
	Settings *SettingsTable
	Alias    *AliasTable
}

type MessageData struct {
	GuildID      string
	ChannelID    string
	MessageID    string
	InsertedTime time.Time
}

func NewDatabase(db *sql.DB) (*Database, error) {
	return &Database{
		pkDb:     db,
		Image:    &ImageTable{pkDb: db},
		Settings: &SettingsTable{pkDb: db},
		Alias:    &AliasTable{pkDb: db},
	}, nil
}

func (db *Database) Initialize() {
	db.pkDb.Exec("CREATE TABLE IF NOT EXISTS image (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT,guildID TEXT, channelID TEXT, messageID TEXT, insertedTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP);")
	db.pkDb.Exec("CREATE TABLE IF NOT EXISTS imageversion (id INTEGER PRIMARY KEY AUTOINCREMENT, imageID INTEGER, url TEXT, FOREIGN KEY(imageID) REFERENCES image(id) ON DELETE CASCADE);")
	db.pkDb.Exec("CREATE TABLE IF NOT EXISTS settings (guildID TEXT PRIMARY KEY, defaultChannel TEXT NOT NULL DEFAULT '', archiveChannel TEXT NOT NULL DEFAULT '')")
	db.pkDb.Exec("CREATE TABLE IF NOT EXISTS alias (id INTEGER PRIMARY KEY AUTOINCREMENT, guildID TEXT NOT NULL DEFAULT '', channelID TEXT NOT NULL DEFAULT '', alias TEXT NOT NULL DEFAULT '')")
}
