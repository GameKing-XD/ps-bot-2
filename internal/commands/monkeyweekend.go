package commands

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/nleeper/goment"
	"github.com/tvanriel/ps-bot-2/internal/bijnaweekend"
)
type MonkeyWeekend struct{}
func NewMonkeyWeekend() *MonkeyWeekend {
        return &MonkeyWeekend{}

}

func (f *MonkeyWeekend) Name() string {
        return "happy funky monkey friday people"
}
func (f *MonkeyWeekend) SkipsPrefix() bool {
        return true
}
func (b *MonkeyWeekend) Apply(ctx *Context) error {
        f, err := os.CreateTemp(os.TempDir(), "bijna-weekend-*.png")
        if err != nil {
                return err
        }
        err = bijnaweekend.MonkeyWeekend(f)
        if err != nil {
                return err
        }
        err = f.Sync()
        if err != nil {
                return err
        }
        _, err = f.Seek(0,0)
        if err != nil {
                return err
        }

        whenweekend, _ := goment.New()
        whenweekend = whenweekend.Local()
	whenweekend.SetDay(5) // friday
	whenweekend.SetHour(16)
	whenweekend.SetMinute(0)
	if whenweekend.IsBefore(goment.New()) {
		whenweekend = whenweekend.Add(7, "days")
	}

        _, err = ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
                Content: fmt.Sprintf("Weekend starts <t:%s:R>", whenweekend.Format("X")),
                Files: []*discordgo.File{
                        {
                                Name: "monkey-weekend.png",
                                Reader: f,
                        },
                },
        })
        return err


}
