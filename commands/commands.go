package commands

import (
	"fmt"
	"strings"

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

func channelIdFromString(source string) (string, error) {
	id := strings.Trim(source, "<>#")
	if len(id) < 1 || len(id) > 20 {
		return "", fmt.Errorf("Invalid text length")
	}
	return id, nil
}
