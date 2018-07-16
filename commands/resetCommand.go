package commands

import (
	"log"

	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type ResetCommand struct {
	pkDb *database.Database
}

func NewResetCommand(pkDb *database.Database) *ResetCommand {
	return &ResetCommand{pkDb}
}

func (c *ResetCommand) Definition() string {
	return "reset"
}

func (c *ResetCommand) Execute(session *discordgo.Session, channel *discordgo.Channel, message *discordgo.MessageCreate) {
	params := parseParameters(c, message.Content)
	defaultCh, _ := c.pkDb.Settings.QueryDefault(channel.GuildID)
	channelParam := func() string {
		if _, ok := params["channel"]; ok {
			return params["channel"]
		}
		return defaultCh
	}()
	if channelParam == "" {
		session.ChannelMessageSend(message.ChannelID, "You need to specify a channel name or \"all\" for all channels.")
		return
	}
	if channelParam == "all" {
		err := c.pkDb.Image.ResetAll(channel.GuildID)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "Cannot reset all page references. Please try again later.")
			log.Println(err)
			return
		}
	} else {
		chID, err := channelIdFromString(channelParam)
		destChannel, err := session.Channel(chID)
		if err != nil || destChannel.GuildID != channel.GuildID {
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
