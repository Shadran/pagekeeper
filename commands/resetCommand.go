package commands

import (
	"log"
	"strings"

	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type ResetCommand struct {
	session *discordgo.Session
	pkDb    *database.Database
}

func NewResetCommand(session *discordgo.Session, pkDb *database.Database) *ResetCommand {
	return &ResetCommand{session, pkDb}
}

func (c *ResetCommand) Definition() string {
	return "reset"
}

func (c *ResetCommand) Execute(channel *discordgo.Channel, message *discordgo.MessageCreate) {
	parts := strings.Split(message.Content, " ")
	if len(parts) < 2 {
		c.session.ChannelMessageSend(message.ChannelID, "You need to specify a channel name or \"all\" for all channels.")
		return
	}
	if parts[1] == "all" {
		err := c.pkDb.Image.ResetAll(channel.GuildID)
		if err != nil {
			c.session.ChannelMessageSend(message.ChannelID, "Cannot reset all page references. Please try again later.")
			log.Println(err)
			return
		}
	} else {
		chID, err := channelIdFromString(parts[1])
		if err != nil {
			c.session.ChannelMessageSend(message.ChannelID, "Invalid channel ID specified.")
			log.Println(err)
			return
		}
		destChannel, err := c.session.Channel(chID)
		if err != nil || destChannel.GuildID != channel.GuildID {
			c.session.ChannelMessageSend(message.ChannelID, "Invalid channel ID specified.")
			return
		}
		c.pkDb.Image.ResetChannel(destChannel.GuildID, destChannel.ID)
		if err != nil {
			c.session.ChannelMessageSend(message.ChannelID, "Cannot reset channel page references. Please try again later.")
			log.Println(err)
			return
		}
	}

	c.session.ChannelMessageSend(message.ChannelID, "References reset successfully.")
}

func (c *ResetCommand) HelpText() string {
	return ""
}
