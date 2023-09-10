package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ready(d *DiscordBot) func(*discordgo.Session, *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("READY")
		for i := range r.Guilds {
			d.repo.LoadGuild(r.Guilds[i].ID)
		}
	}
}

func guildCreate(d *DiscordBot) func(*discordgo.Session, *discordgo.GuildCreate) {
	return func(s *discordgo.Session, gc *discordgo.GuildCreate) {
		d.repo.LoadGuild(gc.ID)
		d.player.Connect(s, gc.ID)
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
