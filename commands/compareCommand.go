package commands

import (
	"bytes"
	"database/sql"
	"log"
	"mime"
	"path"
	"strconv"

	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type CompareCommand struct {
	pkDb *database.Database
}

func NewCompareCommand(pkDb *database.Database) *CompareCommand {
	return &CompareCommand{pkDb}
}

func (c *CompareCommand) Definition() string {
	return "compare"
}

func (c *CompareCommand) Execute(session *discordgo.Session, channel *discordgo.Channel, message *discordgo.MessageCreate) {
	params := parseParameters(c, message.Content)
	defaultCh, _ := c.pkDb.Settings.QueryDefault(channel.GuildID)
	channelParam := func() string {
		if _, ok := params["channel"]; ok {
			return params["channel"]
		}
		return defaultCh
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
		log.Println("Error while getting channel id ", err)
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
	userChannel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Cannot create a direct message channel. Check your settings.")
		return
	}
	session.ChannelMessageSend(message.ChannelID, "Preparing the direct messages, please wait...")
	for i, v := range image.Versions {
		b, err := downloadImage(v.URL)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "There was an error while downloading images.")
			log.Println("Error while downloading images ", err)
			return
		}
		file := &discordgo.File{
			ContentType: mime.TypeByExtension(path.Ext(v.URL)),
			Name:        path.Base(v.URL),
			Reader:      bytes.NewBuffer(b),
		}
		_, err = session.ChannelMessageSendComplex(userChannel.ID, &discordgo.MessageSend{
			Content: image.Title + " version " + strconv.Itoa(i+1),
			File:    file,
		})
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "There was an error while sending the direct message.")
			log.Println("Error while sending direct message ", err)
			return
		}
	}
}

func (c *CompareCommand) HelpText() string {
	return ""
}

func (c *CompareCommand) parameters() []string {
	return []string{"title", "channel"}
}
