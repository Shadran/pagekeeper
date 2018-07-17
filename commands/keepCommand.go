package commands

import (
	"log"
	"mime"
	"path/filepath"
	"strings"

	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type KeepCommand struct {
	baseCommand
}

func NewKeepCommand(pkDb *database.Database, parser *ChannelParser) *KeepCommand {
	return &KeepCommand{newBaseCommand(pkDb, parser)}
}

func (c *KeepCommand) Definition() string {
	return "keep"
}

func (c *KeepCommand) Execute(session *discordgo.Session, channel *discordgo.Channel, message *discordgo.MessageCreate) {
	if len(message.Attachments) == 0 {
		session.ChannelMessageSend(message.ChannelID, "You need to attach an image to the message you want to keep.")
		return
	}
	attachment := message.Attachments[0]
	if !strings.Contains(mime.TypeByExtension(filepath.Ext(attachment.Filename)), "image") {
		session.ChannelMessageSend(message.ChannelID, "The specified attachment is not an image")
		return
	}
	params := parseParameters(c, message.Content)
	destChannel, err := c.parser.Parse(session, channel.GuildID, params["channel"], true, true)
	if err != nil || destChannel == nil || destChannel.GuildID != channel.GuildID {
		session.ChannelMessageSend(message.ChannelID, "Invalid channel ID specified.")
		return
	}
	titleParam := func() string {
		if _, ok := params["title"]; ok {
			return params["title"]
		}
		return strings.TrimSuffix(attachment.Filename, filepath.Ext(attachment.Filename))
	}()

	log.Println("Adding image " + titleParam + " to channel " + destChannel.ID)
	image, err := c.pkDb.Image.Add(titleParam, attachment.URL, database.MessageData{GuildID: destChannel.GuildID, ChannelID: destChannel.ID})
	if image.MessageID != "" {
		log.Println("Image already found, editing")
		_, err = session.ChannelMessageEdit(destChannel.ID, image.MessageID, image.Title+"\r\n"+attachment.URL)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "Cannot add the new image, please try again.")
			return
		}
		currentDefaults, err := c.pkDb.Settings.QueryDefault(channel.GuildID)
		if err != nil {
			log.Println("Error retrieving defaults: ", err)
			return
		}
		if len(image.Versions) > 1 && currentDefaults.ArchiveChannel != "" {
			session.ChannelMessageSend(currentDefaults.ArchiveChannel, "Previous version of image "+image.Title+" from channel <#"+image.ChannelID+">\r\n"+image.Versions[len(image.Versions)-2].URL)
		}
	} else {
		m, err := session.ChannelMessageSend(destChannel.ID, image.Title+"\r\n"+attachment.URL)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "Cannot add the new image, please try again.")
			return
		}
		err = c.pkDb.Image.UpdateLocation(image.ID, database.MessageData{GuildID: destChannel.GuildID, ChannelID: destChannel.ID, MessageID: m.ID})
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "There was an error while adding the new image.")
			log.Println("Error updating references to DB: ", err)
			return
		}
	}
	session.ChannelMessageSend(message.ChannelID, "Page added on channel <#"+destChannel.ID+">")
}

func (c *KeepCommand) HelpText() string {
	return ""
}

func (c *KeepCommand) parameters() []string {
	return []string{"title", "channel"}
}
