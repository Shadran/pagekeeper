package commands

import (
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
	pkDb *database.Database
}

func NewKeepCommand(pkDb *database.Database) *KeepCommand {
	return &KeepCommand{pkDb}
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
	defaultCh, _ := c.pkDb.Settings.QueryDefault(channel.GuildID)
	log.Println("Default channel: ", defaultCh)
	channelParam := func() string {
		if _, ok := params["channel"]; ok {
			log.Println("returning parameter")
			return params["channel"]
		}
		log.Println("returning default")
		return defaultCh
	}()
	if channelParam == "" {
		session.ChannelMessageSend(message.ChannelID, "You need to specify a channel name.")
		return
	}
	chID, err := channelIdFromString(channelParam)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Invalid channel ID specified.")
		log.Println(err)
		return
	}
	destChannel, err := session.Channel(chID)
	if err != nil || destChannel.GuildID != channel.GuildID {
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
