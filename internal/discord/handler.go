package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func ready(d *DiscordBot) func(*discordgo.Session, *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("READY")
		for i := range r.Guilds {
			go d.repo.LoadGuild(r.Guilds[i].ID, r.Guilds[i].Name, r.Guilds[i].IconURL("128"))
		}
	}
}

func guildCreate(d *DiscordBot) func(*discordgo.Session, *discordgo.GuildCreate) {
	return func(s *discordgo.Session, gc *discordgo.GuildCreate) {
		go d.repo.LoadGuild(gc.ID, gc.Name, gc.IconURL("128"))
		go d.player.Connect(s, gc.ID)

	}
}
func messagehandler(d *DiscordBot) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, evt *discordgo.MessageCreate) {
		if evt.Message.WebhookID != "" {
			return
		}

		if evt.Message.Author.Bot {
			return
		}
		d.log.Info("messageCreate",
			zap.String("content", evt.Message.Content),
			zap.String("guild", evt.GuildID),
			zap.String("author", evt.Message.Author.ID),
			zap.String("username", evt.Message.Author.Username),
			zap.String("channel", evt.Message.ChannelID),
		)

		trigger := d.repo.GetPrefix(evt.GuildID)
		if d.exe.HasMatch(trigger, evt.Content) {
			d.exe.Apply(trigger, evt.Message, s)
		}

		for i := range evt.Mentions {
			if evt.Mentions[i].ID == d.conn.State.User.ID {
				d.exe.Apply("<@"+d.conn.State.User.ID+"> ", evt.Message, s)
			}
		}
	}
}
