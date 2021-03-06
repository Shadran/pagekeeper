package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shadran/pagekeeper/commands"
	"github.com/Shadran/pagekeeper/database"
	"github.com/Shadran/pagekeeper/utils"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	SecretToken string
}

func main() {
	manager := utils.NewConfigManager("config.json")
	conf, err := manager.ReadConfig()
	if err != nil {
		log.Fatalln("Cannot read config: ", err)
	}
	session, err := discordgo.New("Bot " + conf.SecretToken)
	if err != nil {
		log.Fatalln("Cannot start discord session: ", err)
	}
	err = session.Open()
	defer session.Close()
	if err != nil {
		log.Fatalln("Error opening discord session: ", err)
	}
	db, err := sql.Open("sqlite3", "pkDatabase.db")
	defer db.Close()
	if err != nil {
		log.Fatalln("Cannot start database: ", err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatalln("Cannot start database: ", err)
	}
	pkDb, err := database.NewDatabase(db)
	if err != nil {
		log.Fatalln("Cannot start database: ", err)
	}
	pkDb.Initialize()
	parser := commands.NewChannelParser(pkDb)

	bot := NewPageKeeper(pkDb)

	bot.Start(session, parser)

	log.Println("Page Keeper is up and running! Press CTRL + C to exit...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func readConfig(path string) (*config, error) {
	result := &config{}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
