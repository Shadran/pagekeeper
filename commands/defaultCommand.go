package commands

import (
	"log"

	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type DefaultCommand struct {
	baseCommand
}

func NewDefaultCommand(pkDb *database.Database, parser *ChannelParser) *DefaultCommand {
	return &DefaultCommand{newBaseCommand(pkDb, parser)}
}

func (c *DefaultCommand) Definition() string {
	return "default"
}

func (c *DefaultCommand) Execute(session *discordgo.Session, channel *discordgo.Channel, message *discordgo.MessageCreate) {
	params := parseParameters(c, message.Content)
	currentDefaults, err := c.pkDb.Settings.QueryDefault(channel.GuildID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Cannot update settings, please try again.")
		log.Println("Error while updating settings: ", err)
		return
	}
	if params["channel"] == "reset" {
		err := c.pkDb.Settings.UpdateDefault(channel.GuildID, "", currentDefaults.ArchiveChannel)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "Cannot update settings, please try again.")
			log.Println("Error while updating settings: ", err)
			return
		}
		session.ChannelMessageSend(message.ChannelID, "Defaults updated successfully.")
		return
	}

	destChannel, err := c.parser.Parse(session, channel.GuildID, params["channel"], false, true)
	if err != nil || destChannel == nil || destChannel.GuildID != channel.GuildID {
		session.ChannelMessageSend(message.ChannelID, "Invalid channel ID specified.")
		return
	}

	err = c.pkDb.Settings.UpdateDefault(destChannel.GuildID, destChannel.ID, currentDefaults.ArchiveChannel)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Cannot update settings, please try again.")
		log.Println("Error while updating settings: ", err)
		return
	}
	session.ChannelMessageSend(message.ChannelID, "Default channel updated successfully.")
}

func (c *DefaultCommand) HelpText() string {
	return ""
}

func (c *DefaultCommand) parameters() []string {
	return []string{"channel"}
}
