package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Definition() string
	Execute(channel *discordgo.Channel, message *discordgo.MessageCreate)
	HelpText() string
}
