package commands

import (
	"log"

	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type DefaultCommand struct {
	pkDb *database.Database
}

func NewDefaultCommand(pkDb *database.Database) *DefaultCommand {
	return &DefaultCommand{pkDb}
}

func (c *DefaultCommand) Definition() string {
	return "default"
}

func (c *DefaultCommand) Execute(session *discordgo.Session, channel *discordgo.Channel, message *discordgo.MessageCreate) {
	params := parseParameters(c, message.Content)
	defaultCh, _ := c.pkDb.Settings.QueryDefault(channel.GuildID)
	channelParam := func() string {
		if _, ok := params["channel"]; ok {
			return params["channel"]
		} else {
			return defaultCh
		}
	}()
	if channelParam == "" {
		session.ChannelMessageSend(message.ChannelID, "You need to specify a channel name.")
		return
	}
	chID, err := channelIdFromString(channelParam)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Invalid channel ID specified.")
		log.Println("Error while getting channel id ", err)
		return
	}
	destChannel, err := session.Channel(chID)
	if err != nil || destChannel.GuildID != channel.GuildID {
		session.ChannelMessageSend(message.ChannelID, "Invalid channel ID specified.")
		return
	}

	err = c.pkDb.Settings.UpdateDefault(destChannel.GuildID, destChannel.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Cannot update settings, please try again.")
		log.Println("Error while updating settings: ", err)
		return
	}
	session.ChannelMessageSend(message.ChannelID, "Defaults updated successfully.")
}

func (c *DefaultCommand) HelpText() string {
	return ""
}

func (c *DefaultCommand) parameters() []string {
	return []string{"channel"}
}
