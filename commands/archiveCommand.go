package commands

import (
	"log"

	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type ArchiveCommand struct {
	baseCommand
}

func NewArchiveCommand(pkDb *database.Database, parser *ChannelParser) *ArchiveCommand {
	return &ArchiveCommand{newBaseCommand(pkDb, parser)}
}

func (c *ArchiveCommand) Definition() string {
	return "archive"
}

func (c *ArchiveCommand) Execute(session *discordgo.Session, channel *discordgo.Channel, message *discordgo.MessageCreate) {
	params := parseParameters(c, message.Content)
	currentDefaults, err := c.pkDb.Settings.QueryDefault(channel.GuildID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Cannot update settings, please try again.")
		log.Println("Error while updating settings: ", err)
		return
	}
	if params["channel"] == "reset" {
		err := c.pkDb.Settings.UpdateDefault(channel.GuildID, currentDefaults.DefaultChannel, "")
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

	err = c.pkDb.Settings.UpdateDefault(destChannel.GuildID, currentDefaults.DefaultChannel, destChannel.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Cannot update settings, please try again.")
		log.Println("Error while updating settings: ", err)
		return
	}
	session.ChannelMessageSend(message.ChannelID, "Archive channel updated successfully.")
}

func (c *ArchiveCommand) HelpText() string {
	return ""
}

func (c *ArchiveCommand) parameters() []string {
	return []string{"channel"}
}
