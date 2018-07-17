package database

import (
	"database/sql"
	"log"
	"time"
)

type ImageTable struct {
	pkDb *sql.DB
}

type Image struct {
	ID       int64
	Title    string
	Versions []*ImageVersion
	MessageData
}

type ImageVersion struct {
	ID  int64
	URL string
}

func (t *ImageTable) Add(title string, url string, source MessageData) (*Image, error) {
	tx, err := t.pkDb.Begin()
	if err != nil {
		return nil, err
	}
	image, err := t.QueryByTitleAndLocation(title, source.GuildID, source.ChannelID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		image, err = t.addNewImage(title, source, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	newVersion, err := t.addVersion(image.ID, url, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	image.Versions = append(image.Versions, newVersion)
	return image, nil
}

func (t *ImageTable) QueryByTitleAndLocation(title string, guildID string, channelID string) (*Image, error) {
	stmt, err := t.pkDb.Prepare(`SELECT image.id, image.title, image.guildID, image.channelID, image.messageID, image.insertedTime
								 FROM image 
								 WHERE image.title = ? AND image.guildID = ? AND image.channelID = ?`)
	if err != nil {
		return nil, err
	}
	image := &Image{}
	log.Println("Searching images with ", title, guildID, channelID)
	if err := stmt.QueryRow(title, guildID, channelID).Scan(&image.ID, &image.Title, &image.GuildID, &image.ChannelID, &image.MessageID, &image.InsertedTime); err != nil {
		log.Println("Error while searching ", err)
		return nil, err
	}
	image.Versions, err = t.getVersions(image.ID)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (t *ImageTable) QueryOrdered(guildID string, channelID string) ([]*Image, error) {
	stmt, err := t.pkDb.Prepare(`SELECT image.id, image.title, image.guildID, image.channelID, image.messageID, image.insertedTime
								 FROM image 
								 WHERE image.guildID = ? AND image.channelID = ?
								 ORDER BY  CAST(title AS INTEGER), image.title`)
	if err != nil {
		return nil, err
	}
	images := []*Image{}
	rows, err := stmt.Query(guildID, channelID)
	for rows.Next() {
		image := &Image{}
		err = rows.Scan(&image.ID, &image.Title, &image.GuildID, &image.ChannelID, &image.MessageID, &image.InsertedTime)
		if err != nil {
			log.Println("Error while retrieving messages: ", err)
			return nil, err
		}
		image.Versions, err = t.getVersions(image.ID)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}

func (t *ImageTable) UpdateLocation(id int64, source MessageData) error {
	tx, err := t.pkDb.Begin()
	if err != nil {
		return err
	}
	if source.InsertedTime == (time.Time{}) {
		stmt, err := t.pkDb.Prepare(`UPDATE image SET guildID = ?, channelID = ?, messageID = ? WHERE id = ?`)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = stmt.Exec(source.GuildID, source.ChannelID, source.MessageID, id)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		stmt, err := t.pkDb.Prepare(`UPDATE image SET guildID = ?, channelID = ?, messageID = ?, insertedTime = ? WHERE id = ?`)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = stmt.Exec(source.GuildID, source.ChannelID, source.MessageID, source.InsertedTime, id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (t *ImageTable) ResetChannel(guildID string, channelID string) error {
	tx, err := t.pkDb.Begin()
	if err != nil {
		return err
	}
	stmt, err := t.pkDb.Prepare(`DELETE FROM image WHERE guildID = ? AND channelID = ?`)
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

func (t *ImageTable) ResetAll(guildID string) error {
	tx, err := t.pkDb.Begin()
	if err != nil {
		return err
	}
	stmt, err := t.pkDb.Prepare(`DELETE FROM image WHERE guildID = ?`)
	if err != nil {
		log.Println("Error while resetting: ", err)
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(guildID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (t *ImageTable) RemoveImage(guildID string, channelID string, imageID int64) error {
	tx, err := t.pkDb.Begin()
	if err != nil {
		return err
	}
	stmt, err := t.pkDb.Prepare(`DELETE FROM image WHERE guildID = ? AND channelID = ? AND id = ?`)
	if err != nil {
		log.Println("Error while removing image: ", err)
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(guildID, channelID, imageID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (t *ImageTable) addNewImage(title string, source MessageData, tx *sql.Tx) (*Image, error) {
	image := &Image{Title: title, MessageData: MessageData{source.GuildID, source.ChannelID, source.MessageID, source.InsertedTime}}
	stmt, err := tx.Prepare("INSERT INTO image (title, guildID, channelID, messageID) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("Error while adding new image: ", err)
		return nil, err
	}
	result, err := stmt.Exec(title, source.GuildID, source.ChannelID, source.MessageID)
	if err != nil {
		log.Println("Error while adding new image: ", err)
		return nil, err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error while adding new image: ", err)
		return nil, err
	}
	image.ID = lastID
	return image, nil
}

func (t *ImageTable) getVersions(imageID int64) ([]*ImageVersion, error) {
	stmt, err := t.pkDb.Prepare(`SELECT id, url
								 FROM imageversion WHERE imageID = ?`)
	if err != nil {
		log.Println("Error while retrieving versions: ", err)
		return nil, err
	}
	versions := []*ImageVersion{}
	rows, err := stmt.Query(imageID)
	for rows.Next() {
		version := &ImageVersion{}
		err = rows.Scan(&version.ID, &version.URL)
		if err != nil {
			log.Println("Error while retrieving versions: ", err)
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

func (t *ImageTable) addVersion(imageID int64, url string, tx *sql.Tx) (*ImageVersion, error) {
	version := &ImageVersion{URL: url}
	stmt, err := tx.Prepare("INSERT INTO imageversion (imageID, url) VALUES (?, ?)")
	if err != nil {
		log.Println("Error while adding versions: ", err)
		return nil, err
	}
	result, err := stmt.Exec(imageID, url)
	if err != nil {
		log.Println("Error while adding versions: ", err)
		return nil, err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error while adding versions: ", err)
		return nil, err
	}
	version.ID = lastID
	return version, nil
}
