package commands

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/nleeper/goment"
	"github.com/tvanriel/ps-bot-2/internal/bijnaweekend"
)
type FuehrerWeekend struct{}
func NewFuehrerWeekend() *FuehrerWeekend {
        return &FuehrerWeekend{}

}

func (f *FuehrerWeekend) Name() string {
        return "Meine herren, es ist fast wochenende"
}
func (f *FuehrerWeekend) SkipsPrefix() bool {
        return true
}
func (b *FuehrerWeekend) Apply(ctx *Context) error {
        f, err := os.CreateTemp(os.TempDir(), "bijna-weekend-*.png")
        if err != nil {
                return err
        }
        err = bijnaweekend.FuehrerWeekend(f)
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
                                Name: "fuehrer-weekend.png",
                                Reader: f,
                        },
                },
        })
        return err


}
