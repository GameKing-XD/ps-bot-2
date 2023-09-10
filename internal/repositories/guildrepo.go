package repositories

import (
	"errors"

	"gorm.io/gorm"
)

const DEFAULT_PREFIX = "!"

type GuildModel struct {
	gorm.Model
	Snowflake    string
	Voicechannel string
	Prefix       string
}

type GuildRepository struct {
	db *gorm.DB
}

func NewGuildRepository(mysql *gorm.DB) *GuildRepository {
	return &GuildRepository{
		db: mysql,
	}
}

func MigrateGuildRepo(g *GuildRepository) error {
	return g.db.AutoMigrate(
		&GuildModel{},
	)
}

func (g *GuildRepository) LoadGuild(guildId string) {
	var model GuildModel
	tx := g.db.First(&model, "snowflake = ?", guildId)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		g.db.Create(&GuildModel{
			Snowflake: guildId,
			Prefix:    DEFAULT_PREFIX,
		})
	}

}

func (g *GuildRepository) GetVoiceChannels() map[string]string {

	m := make(map[string]string)

	var models []GuildModel
	tx := g.db.Find(models, "voicechannel != ''")
	if tx.Error != nil {
		return m
	}

	for i := range models {
		m[models[i].Snowflake] = models[i].Voicechannel
	}
	return m

}

func (g *GuildRepository) UpdatePrefix(guildId, prefix string) {
	g.db.Model(&GuildModel{}).Where("snowflake", guildId).Update("prefix", prefix)
}

func (g *GuildRepository) GetPrefix(guildId string) string {
	var m GuildModel
	tx := g.db.First(&m, "snowflake = ?", guildId)
	if tx.Error != nil {
		return DEFAULT_PREFIX
	}
	if m.Snowflake == "" {
		return DEFAULT_PREFIX
	}
	return m.Prefix

}
func (g *GuildRepository) JoinVoiceChannel(guildId, channelId string) {
	g.db.Model(&GuildModel{}).Where("snowflake = ?", guildId).UpdateColumn("voicechannel", channelId)

}

func (g *GuildRepository) GetVoiceChannel(guildId string) string {
	var channelId string

	tx := g.db.Model(&GuildModel{}).Where("snowflake = ?", guildId).Pluck("voicechannel", &channelId)

	if tx.Error != nil {
		return ""
	}

	return channelId

}
