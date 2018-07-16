package commands

import (
	"log"
	"sort"

	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type OrderCommand struct {
	pkDb *database.Database
}

func NewOrderCommand(pkDb *database.Database) *OrderCommand {
	return &OrderCommand{pkDb}
}

func (c *OrderCommand) Definition() string {
	return "order"
}

func (c *OrderCommand) Execute(session *discordgo.Session, channel *discordgo.Channel, message *discordgo.MessageCreate) {
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
	titleOrderedImages, err := c.pkDb.Image.QueryOrdered(destChannel.GuildID, destChannel.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Error while ordering images.")
		return
	}
	timeOrderedImages := make([]*database.Image, len(titleOrderedImages))
	copy(timeOrderedImages, titleOrderedImages)
	sort.Slice(timeOrderedImages, func(i, j int) bool {
		return timeOrderedImages[i].InsertedTime.Before(timeOrderedImages[j].InsertedTime)
	})
	for i := range titleOrderedImages {
		log.Println("Moving image ", titleOrderedImages[i].Title, " to message sent at ", timeOrderedImages[i].InsertedTime.String())
		_, err = session.ChannelMessageEdit(destChannel.ID, timeOrderedImages[i].MessageID, titleOrderedImages[i].Title+"\r\n"+titleOrderedImages[i].Versions[len(titleOrderedImages[i].Versions)-1].URL)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "Error while ordering images")
			log.Println("Cannot edit message: ", err)
			return
		}
		c.pkDb.Image.UpdateLocation(titleOrderedImages[i].ID, database.MessageData{GuildID: destChannel.GuildID, ChannelID: destChannel.ID, MessageID: timeOrderedImages[i].MessageID, InsertedTime: timeOrderedImages[i].InsertedTime})
	}
	session.ChannelMessageSend(message.ChannelID, "Images ordered by title successfully.")
}

func (c *OrderCommand) HelpText() string {
	return ""
}

func (c *OrderCommand) parameters() []string {
	return []string{"channel"}
}
