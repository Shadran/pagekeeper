package database

import (
	"database/sql"
)

type Database struct {
	pkDb  *sql.DB
	Image *ImageTable
}

type MessageData struct {
	GuildID   string
	ChannelID string
	MessageID string
}

func NewDatabase(db *sql.DB) (*Database, error) {
	return &Database{pkDb: db, Image: &ImageTable{pkDb: db}}, nil
}

func (db *Database) Initialize() {
	db.pkDb.Exec("CREATE TABLE IF NOT EXISTS image (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, guildID TEXT, channelID TEXT, messageID TEXT);")
	db.pkDb.Exec("CREATE TABLE IF NOT EXISTS imageversion (id INTEGER PRIMARY KEY AUTOINCREMENT, imageID INTEGER, url TEXT, FOREIGN KEY(imageID) REFERENCES image(id) ON DELETE CASCADE);")
}
