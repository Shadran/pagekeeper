package commands

import (
	"log"

	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type ResetCommand struct {
	baseCommand
}

func NewResetCommand(pkDb *database.Database, parser *ChannelParser) *ResetCommand {
	return &ResetCommand{newBaseCommand(pkDb, parser)}
}

func (c *ResetCommand) Definition() string {
	return "reset"
}

func (c *ResetCommand) Execute(session *discordgo.Session, channel *discordgo.Channel, message *discordgo.MessageCreate) {
	params := parseParameters(c, message.Content)
	if params["channel"] == "all" {
		err := c.pkDb.Image.ResetAll(channel.GuildID)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "Cannot reset all page references. Please try again later.")
			log.Println(err)
			return
		}
	} else {
		destChannel, err := c.parser.Parse(session, channel.GuildID, params["channel"], true, true)
		if err != nil || destChannel == nil || destChannel.GuildID != channel.GuildID {
			session.ChannelMessageSend(message.ChannelID, "Invalid channel ID specified.")
			return
		}
		c.pkDb.Image.ResetChannel(destChannel.GuildID, destChannel.ID)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "Cannot reset channel page references. Please try again later.")
			log.Println(err)
			return
		}
	}

	session.ChannelMessageSend(message.ChannelID, "References reset successfully.")
}

func (c *ResetCommand) HelpText() string {
	return ""
}

func (c *ResetCommand) parameters() []string {
	return []string{"channel"}
}
