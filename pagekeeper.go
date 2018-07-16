package main

import (
	"strings"

	"github.com/Shadran/pagekeeper/database"

	"github.com/Shadran/pagekeeper/commands"
	"github.com/bwmarrin/discordgo"
)

const prefix string = "pk:"

type PageKeeper struct {
	pkDb     *database.Database
	commands map[string]commands.Command
}

func NewPageKeeper(db *database.Database) *PageKeeper {
	return &PageKeeper{pkDb: db, commands: map[string]commands.Command{}}
}

func (pk *PageKeeper) Start(session *discordgo.Session) {
	list := []commands.Command{
		commands.NewKeepCommand(pk.pkDb),
		commands.NewResetCommand(pk.pkDb),
		commands.NewOrderCommand(pk.pkDb),
		commands.NewDefaultCommand(pk.pkDb),
		commands.NewRemoveCommand(pk.pkDb),
		commands.NewCompareCommand(pk.pkDb),
	}

	for _, c := range list {
		pk.commands[c.Definition()] = c
	}

	session.AddHandler(pk.messageCreate)
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
	c.Execute(session, channel, message)
}
