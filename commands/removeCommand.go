package commands

import (
	"database/sql"
	"log"

	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type RemoveCommand struct {
	baseCommand
}

func NewRemoveCommand(pkDb *database.Database, parser *ChannelParser) *RemoveCommand {
	return &RemoveCommand{newBaseCommand(pkDb, parser)}
}

func (c *RemoveCommand) Definition() string {
	return "remove"
}

func (c *RemoveCommand) Execute(session *discordgo.Session, channel *discordgo.Channel, message *discordgo.MessageCreate) {
	params := parseParameters(c, message.Content)
	destChannel, err := c.parser.Parse(session, channel.GuildID, params["channel"], true, true)
	if err != nil || destChannel == nil || destChannel.GuildID != channel.GuildID {
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
