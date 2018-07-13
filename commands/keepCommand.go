package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type KeepCommand struct {
	session *discordgo.Session
	pkDb    *database.Database
}

func NewKeepCommand(session *discordgo.Session, pkDb *database.Database) *KeepCommand {
	return &KeepCommand{session, pkDb}
}

func (c *KeepCommand) Definition() string {
	return "keep"
}

func (c *KeepCommand) Execute(channel *discordgo.Channel, message *discordgo.MessageCreate) {
	if len(message.Attachments) == 0 {
		c.session.ChannelMessageSend(message.ChannelID, "You need to attach an image to the message you want to keep.")
		return
	}
	attachment := message.Attachments[0]
	if !strings.Contains(mime.TypeByExtension(filepath.Ext(attachment.Filename)), "image") {
		c.session.ChannelMessageSend(message.ChannelID, "The specified attachment is not an image")
		return
	}
	parts := strings.Split(message.Content, " ")
	if len(parts) < 3 {
		c.session.ChannelMessageSend(message.ChannelID, "You need to specify a title and a channel for the page.")
		return
	}
	chID, err := channelIdFromString(parts[2])
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
	log.Println("Adding image " + parts[1] + " to channel " + destChannel.ID)
	image, err := c.pkDb.Image.Add(parts[1], attachment.URL, database.MessageData{GuildID: destChannel.GuildID, ChannelID: destChannel.ID})
	if image.MessageID != "" {
		log.Println("Image already found, editing")
		_, err = c.session.ChannelMessageEdit(destChannel.ID, image.MessageID, image.Title+"\r\n"+attachment.URL)
		if err != nil {
			c.session.ChannelMessageSend(message.ChannelID, "Cannot add the new image, please try again.")
			return
		}
	} else {
		m, err := c.session.ChannelMessageSend(destChannel.ID, image.Title+"\r\n"+attachment.URL)
		if err != nil {
			c.session.ChannelMessageSend(message.ChannelID, "Cannot add the new image, please try again.")
			return
		}
		c.pkDb.Image.UpdateLocation(image.ID, database.MessageData{GuildID: destChannel.GuildID, ChannelID: destChannel.ID, MessageID: m.ID})
	}
	c.session.ChannelMessageSend(message.ChannelID, "Page added on channel "+parts[2])
}

func (c *KeepCommand) HelpText() string {
	return ""
}

func downloadImage(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func channelIdFromString(source string) (string, error) {
	if len(source) < 2 || len(source) > 23 {
		return "", fmt.Errorf("Invalid source text length")
	}
	id := strings.Trim(source, "<>#")
	return id, nil
}
