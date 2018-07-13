package main

import (
	"strings"

	"github.com/Shadran/pagekeeper/database"

	"github.com/Shadran/pagekeeper/commands"
	"github.com/bwmarrin/discordgo"
)

const prefix string = "pk:"

type PageKeeper struct {
	session  *discordgo.Session
	pkDb     *database.Database
	commands map[string]commands.Command
}

func NewPageKeeper(session *discordgo.Session, db *database.Database) *PageKeeper {
	return &PageKeeper{session: session, pkDb: db, commands: map[string]commands.Command{}}
}

func (pk *PageKeeper) Start() {
	list := []commands.Command{
		commands.NewKeepCommand(pk.session, pk.pkDb),
		commands.NewResetCommand(pk.session, pk.pkDb),
	}

	for _, c := range list {
		pk.commands[c.Definition()] = c
	}

	pk.session.AddHandler(pk.messageCreate)
}

func (pk *PageKeeper) messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if !strings.HasPrefix(message.Content, prefix) {
		return
	}
	c, ok := pk.commands[strings.Split(strings.Split(message.Content, " ")[0], ":")[1]]
	if !ok {
		return
	}
	channel, _ := session.Channel(message.ChannelID)
	c.Execute(channel, message)
}
