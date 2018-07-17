package commands

import (
	"github.com/Shadran/pagekeeper/database"
	"github.com/bwmarrin/discordgo"
)

type AliasCommand struct {
	baseCommand
}

func NewAliasCommand(pkDb *database.Database, parser *ChannelParser) *AliasCommand {
	return &AliasCommand{newBaseCommand(pkDb, parser)}
}

func (c *AliasCommand) Definition() string {
	return "alias"
}

func (c *AliasCommand) Execute(session *discordgo.Session, channel *discordgo.Channel, message *discordgo.MessageCreate) {
	params := parseParameters(c, message.Content)
	if params["alias"] == "" {
		session.ChannelMessageSend(message.ChannelID, "You need to specify an alias for the channel.")
		return
	}
	destChannel, err := c.parser.Parse(session, channel.GuildID, params["channel"], false, false)
	if err != nil || destChannel == nil || destChannel.GuildID != channel.GuildID {
		session.ChannelMessageSend(message.ChannelID, "Invalid channel ID specified.")
		return
	}

	err = c.parser.EditAlias(channel.GuildID, destChannel.ID, params["alias"])
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Error while updating aliases.")
		return
	}

	session.ChannelMessageSend(message.ChannelID, "Alias set successfully.")
}

func (c *AliasCommand) HelpText() string {
	return ""
}

func (c *AliasCommand) parameters() []string {
	return []string{"channel", "alias"}
}
