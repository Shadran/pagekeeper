package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Shadran/pagekeeper/database"
)

type ChannelParser struct {
	pkDb    *database.Database
	aliases map[string]map[string]string
}

func NewChannelParser(pkDb *database.Database) *ChannelParser {
	return &ChannelParser{pkDb, map[string]map[string]string{}}
}

func (p *ChannelParser) Parse(session *discordgo.Session, guildId string, channelString string, considerDefaults bool, considerAliases bool) (*discordgo.Channel, error) {
	currentDefaults, err := p.pkDb.Settings.QueryDefault(guildId)
	if err != nil {
		log.Println("Cannot read current settings ", err)
		return nil, err
	}
	channelID := func() string {
		if channelString == "" {
			if considerDefaults {
				return currentDefaults.DefaultChannel
			}
			return ""
		}
		if cm, ok := p.aliases[guildId]; ok && considerAliases {
			if c, ok := cm[channelString]; ok {
				return c
			}
		}
		c, err := channelIdFromString(channelString)
		if err != nil {
			return ""
		}
		return c
	}()
	if channelID == "" {
		return nil, nil
	}
	destChannel, err := session.Channel(channelID)
	if err != nil || destChannel.GuildID != guildId {
		return nil, err
	}
	return destChannel, nil
}

func (p *ChannelParser) EditAlias(guildID string, channelID string, alias string) error {
	err := p.pkDb.Alias.InsertUpdateAlias(guildID, channelID, alias)
	if err != nil {
		return err
	}
	if _, ok := p.aliases[guildID]; !ok {
		p.aliases[guildID] = map[string]string{}
	}
	p.aliases[guildID][alias] = channelID
	return nil
}

func (p *ChannelParser) ReloadAliases() error {
	aliases, err := p.pkDb.Alias.QueryAliases()
	if err != nil {
		return err
	}
	p.aliases = map[string]map[string]string{}
	for _, a := range aliases {
		if _, ok := p.aliases[a.GuildID]; !ok {
			p.aliases[a.GuildID] = map[string]string{}
		}
		p.aliases[a.GuildID][a.Alias] = a.ChannelID
	}
	return nil
}

func channelIdFromString(source string) (string, error) {
	id := strings.Trim(source, "<>#")
	if len(id) < 1 || len(id) > 20 {
		return "", fmt.Errorf("Invalid text length")
	}
	return id, nil
}
