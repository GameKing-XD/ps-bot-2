package bijnaweekend

import (
	"image"
	"image/color"
	"io"
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/muesli/gamut"
)

const (
	bijnaWeekendTextWidth    = 241
	bijnaWeekendTextHeight   = 30
	bijnaWeekendMargin       = 15
	bijnaWeekendDoubleMargin = 2 * bijnaWeekendMargin
	elstakwidth              = 357
	elstakheight             = 480
)

type Palette struct {
	Background1 color.Color
	Background2 color.Color
	Text        color.Color
	TextShadow  color.Color
}

func MonkeyWeekend(w io.Writer) error {
	e := getRandomMonkey()
	p, err := getPalette()
	if err != nil {
		return err
	}
	im := elstakImage(e)
	dc := gg.NewContext(
		im.Bounds().Dx(),
		im.Bounds().Dy(),
	)

	drawBackground(dc, p)
	drawElstak(im, dc)
	drawText(e, p, dc)

	return dc.EncodePNG(w)
}

func FuehrerWeekend(w io.Writer) error {
	e := getRandomFuehrer()
	p, err := getPalette()
	if err != nil {
		return err
	}
	im := elstakImage(e)
	dc := gg.NewContext(
		im.Bounds().Dx(),
		im.Bounds().Dy(),
	)

	drawBackground(dc, p)
	drawElstak(im, dc)
	drawText(e, p, dc)

	return dc.EncodePNG(w)
}

func BijnaWeekend(w io.Writer) error {

	e := getRandomElstak()
	p, err := getPalette()
	if err != nil {
		return err
	}
	im := elstakImage(e)
	dc := gg.NewContext(
		im.Bounds().Dx(),
		im.Bounds().Dy(),
	)

	drawBackground(dc, p)
	drawElstak(im, dc)
	drawText(e, p, dc)

	return dc.EncodePNG(w)

}

func getPalette() (Palette, error) {
	p, err := gamut.Generate(4, gamut.PastelGenerator{})
	if err != nil {
		return Palette{}, err
	}
	return Palette{
		Background1: gamut.Lighter(p[0], 0.1),
		Background2: gamut.Lighter(p[1], 0.1),
		Text:        gamut.Darker(p[3], 0.3),
		TextShadow:  gamut.Darker(p[3], 0.7),
	}, nil
}

func drawBackground(dc *gg.Context, palette Palette) {
	bgPick := rand.Intn(3)

	var bg gg.Gradient
	switch bgPick {
	case 0:
		bg = gg.NewLinearGradient(0, 0, 0, float64(dc.Width()))
	case 1:
		bg = gg.NewLinearGradient(0, 0, float64(dc.Height()), float64(dc.Width()))
	case 2:
		bg = gg.NewLinearGradient(0, float64(dc.Width()), 0, 0)

	}

	bg.AddColorStop(0, palette.Background1)
	bg.AddColorStop(1, palette.Background2)
	dc.SetFillStyle(bg)
	dc.DrawRectangle(0, 0, float64(dc.Width()), float64(dc.Height()))
	dc.Fill()
}

type ElstakConfig struct {
	Filename string
	TextX    float64
	TextY    float64
	Font     string
}

var fuehrers = []ElstakConfig{
	{
		Filename: "assets/hitler1.png",
		TextX:    15,
		TextY:    100,
		Font:     "assets/fraktur.ttf",
	},
	{
		Filename: "assets/hitler2.png",
		TextX:    30,
		TextY:    300,
		Font:     "assets/fraktur.ttf",
	},
	{
		Filename: "assets/hitler3.png",
		TextX:    30,
		TextY:    300,
		Font:     "assets/fraktur.ttf",
	},
	{
		Filename: "assets/hitler4.png",
		TextX:    30,
		TextY:    300,
		Font:     "assets/fraktur.ttf",
	},
}

var elstaks = []ElstakConfig{
	{
		Filename: "assets/elstak-left.png",
		TextX:    15,
		TextY:    100,
		Font:     "assets/impact.ttf",
	},
	{
		Filename: "assets/elstak-up.png",
		TextX:    30,
		TextY:    50,
		Font:     "assets/impact.ttf",
	},
	{
		Filename: "assets/elstak-yay.png",
		TextX:    30,
		TextY:    50,
		Font:     "assets/impact.ttf",
	},
}

var monkeys = []ElstakConfig{
	{
		Filename: "assets/monkey-1.png",
		TextX:    15,
		TextY:    100,
		Font:     "assets/impact.ttf",
	},
}

func getRandomFuehrer() ElstakConfig {
	return fuehrers[rand.Intn(len(fuehrers))]
}

func getRandomElstak() ElstakConfig {
	return elstaks[rand.Intn(len(elstaks))]
}

func getRandomMonkey() ElstakConfig {
	return monkeys[rand.Intn(len(monkeys))]
}
func elstakImage(elstakConfig ElstakConfig) image.Image {
	im, _ := gg.LoadImage(elstakConfig.Filename)
	return im
}

func drawElstak(im image.Image, dc *gg.Context) {
	dc.MoveTo(0, 0)
	dc.DrawImage(im, 0, 0)
}
func drawText(elstakConfig ElstakConfig, palette Palette, dc *gg.Context) {
	dc.LoadFontFace(elstakConfig.Font, 30)

	dc.SetFillStyle(gg.NewSolidPattern(palette.TextShadow))

	distance := float64(2)
	dc.DrawStringAnchored("Bijna Weekend", elstakConfig.TextX+distance, elstakConfig.TextY+distance, 0, 0)
	dc.SetFillStyle(gg.NewSolidPattern(palette.Text))

	dc.DrawStringAnchored("Bijna Weekend", elstakConfig.TextX, elstakConfig.TextY, 0, 0)
}
