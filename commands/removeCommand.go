package commands

import (
	"database/sql"
	"log"

	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type RemoveCommand struct {
	pkDb *database.Database
}

func NewRemoveCommand(pkDb *database.Database) *RemoveCommand {
	return &RemoveCommand{pkDb}
}

func (c *RemoveCommand) Definition() string {
	return "remove"
}

func (c *RemoveCommand) Execute(session *discordgo.Session, channel *discordgo.Channel, message *discordgo.MessageCreate) {
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
	if _, ok := params["title"]; !ok {
		session.ChannelMessageSend(message.ChannelID, "You need to specify an image title.")
		return
	}
	chID, err := channelIdFromString(channelParam)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Invalid channel ID specified.")
		log.Println("Error during channel request: ", err)
		return
	}
	destChannel, err := session.Channel(chID)
	if err != nil || destChannel.GuildID != channel.GuildID {
		session.ChannelMessageSend(message.ChannelID, "Invalid channel ID specified.")
		return
	}

	image, err := c.pkDb.Image.QueryByTitleAndLocation(params["title"], destChannel.GuildID, destChannel.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			session.ChannelMessageSend(message.ChannelID, "There was an error retrieving the specified image.")
			log.Println("Error during image query: ", err)
			return
		}
		session.ChannelMessageSend(message.ChannelID, "Cannot find an image with title "+params["title"]+".")
		return
	}

	err = session.ChannelMessageDelete(image.ChannelID, image.MessageID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Cannot delete the specified image, please try again.")
		log.Println("Error during message removal: ", err)
		return
	}
	err = c.pkDb.Image.RemoveImage(image.GuildID, image.ChannelID, image.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Cannot delete the specified image, please try again.")
		log.Println("Error during image db removal: ", err)
		return
	}
	session.ChannelMessageSend(message.ChannelID, "Image deleted successfully.")
}

func (c *RemoveCommand) HelpText() string {
	return ""
}

func (c *RemoveCommand) parameters() []string {
	return []string{"title", "channel"}
}
