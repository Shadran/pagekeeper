package commands

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Definition() string
	Execute(session *discordgo.Session, channel *discordgo.Channel, message *discordgo.MessageCreate)
	HelpText() string
}

type parametricCommand interface {
	parameters() []string
}

type baseCommand struct {
	pkDb   *database.Database
	parser *ChannelParser
}

func newBaseCommand(pkDb *database.Database, parser *ChannelParser) baseCommand {
	return baseCommand{pkDb, parser}
}

func parseParameters(c parametricCommand, messageText string) map[string]string {
	parts := strings.Split(messageText, " ")[1:]
	parameters := c.parameters()
	result := map[string]string{}
	for i, p := range parameters {
		if i < len(parts) {
			result[p] = parts[i]
		}
	}
	return result
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
