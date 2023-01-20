package tabdrawer

import (
	"fmt"
	"image"
	"image/color"
	"sort"
	"strings"

	"github.com/Tnze/go-mc/chat"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/google/uuid"
	"golang.org/x/image/font"
)

var DefaultChatColorCodes = map[string]color.RGBA{
	"black":        {0, 0, 0, 255},
	"dark_blue":    {0, 0, 170, 255},
	"dark_green":   {0, 170, 0, 255},
	"dark_aqua":    {0, 170, 170, 255},
	"dark_red":     {170, 0, 0, 255},
	"dark_purple":  {170, 0, 170, 255},
	"gold":         {255, 170, 0, 255},
	"gray":         {170, 170, 170, 255},
	"dark_gray":    {85, 85, 85, 255},
	"blue":         {85, 85, 255, 255},
	"green":        {85, 255, 85, 255},
	"aqua":         {85, 255, 255, 255},
	"red":          {255, 85, 85, 255},
	"light_purple": {255, 85, 255, 255},
	"yellow":       {255, 255, 85, 255},
	"white":        {255, 255, 255, 255},
}

func DefaultLatencyColoring(ping int) color.Color {
	if ping < 60 {
		return color.RGBA{0, 255, 0, 255}
	} else if ping < 120 {
		return color.RGBA{105, 155, 0, 255}
	} else if ping < 240 {
		return color.RGBA{180, 90, 0, 255}
	} else if ping < 600 {
		return color.RGBA{255, 60, 60, 255}
	} else {
		return color.RGBA{255, 60, 60, 255}
	}
}

var (
	DefaultDiscordBackgroundColor = color.RGBA{R: 0x36, G: 0x39, B: 0x3f, A: 0xff}
	DefaultPlayerBackgroundColor  = color.RGBA{0, 0, 0, 150}
)

type LatencyStyling int

const (
	LatencyNumber LatencyStyling = iota
	LatencyNumberMs
)

type TabParameters struct {
	LatencyColoring       func(int) color.Color
	LatencyStyle          LatencyStyling
	ChatColorCodes        map[string]color.RGBA
	BackgroundColor       color.Color
	PlayerBackgroundColor color.Color
	RowSpacing            float64
	RowAdditionalHeight   float64
	ColumnSpacing         float64
	MaxRows               int
	FontColor             color.Color
	Font                  font.Face
	LineSpacing           float64
	DebugTopBottom        bool
	DebugHeight           bool
}

type TabPlayer struct {
	Name        string
	Ping        int
	HeadTexture image.Image
	Gamemode    string
}

func DrawTab(players map[uuid.UUID]TabPlayer, tabtop, tabbottom *chat.Message, params TabParameters) image.Image {

	mctx := gg.NewContext(500, 500) // measuring context
	mctx.SetFontFace(params.Font)

	cols := len(players) / params.MaxRows
	if len(players)%params.MaxRows != 0 {
		cols++
	}

	keys := make([]uuid.UUID, 0, len(players))
	for u := range players {
		keys = append(keys, u)
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(players[keys[i]].Name, players[keys[j]].Name) < 0
	})

	pmw, pmh := float64(0), float64(0)
	for _, v := range players {
		w, h := mctx.MeasureString(fmt.Sprint(v.Name, v.Ping, "    ms"))
		if pmw < w {
			pmw = w
		}
		if pmh < h {
			pmh = h
		}
	}
	tabw := float64(float64(cols)*(pmw+pmh+params.ColumnSpacing) + 16)
	tabtopw, tabtoph := measureMaxLine(mctx, *tabtop, params.LineSpacing)
	tabbottomw, tabbottomh := measureMaxLine(mctx, *tabbottom, params.LineSpacing)
	_, lineh := mctx.MeasureString(" ")
	if tabw < tabtopw {
		tabw = tabtopw + 16
	}
	if tabw < tabbottomw {
		tabw = tabbottomw + 16
	}

	colw := pmw + pmh
	rowh := pmh + params.RowAdditionalHeight*2

	tabh := tabtoph + lineh + (rowh+params.RowSpacing)*(float64(params.MaxRows)) + lineh + tabbottomh + lineh
	c := gg.NewContext(int(tabw), int(tabh))
	c.SetColor(params.BackgroundColor)
	c.Clear()
	c.SetFontFace(params.Font)

	if params.DebugHeight {
		c.SetColor(color.RGBA{255, 0, 0, 255})
		c.DrawRectangle(tabw/2-8, 1, 16, tabtoph+lineh)
		c.Fill()
		c.DrawRectangle(tabw/2+8, 1+tabtoph+lineh, 16, ((rowh+params.RowSpacing)*(float64(params.MaxRows)) + lineh))
		c.Fill()
	}

	topf := fragmentMessage(c, gg.AlignCenter, *tabtop, tabw/2, lineh, params.ChatColorCodes)
	topmy := float64(0)
	for _, f := range topf {
		c.SetColor(f.color)
		c.DrawString(f.str, f.x, f.y)
		if topmy < f.y {
			topmy = f.y
		}
		if params.DebugTopBottom {
			w, h := c.MeasureString(f.str)
			c.SetColor(color.RGBA{255, 0, 0, 255})
			c.DrawRectangle(f.x, f.y-h, w, h)
			c.Stroke()
		}
	}

	plc := 0
	for col := 0; col < cols; col++ {
		for row := 0; row < params.MaxRows; row++ {
			if plc > len(keys)-1 {
				break
			}
			pl := players[keys[plc]]
			c.SetColor(params.PlayerBackgroundColor)
			rowx := tabw/2 - (float64(cols)*(colw+params.ColumnSpacing))/2 + float64(col)*(colw+params.ColumnSpacing) + params.ColumnSpacing/2
			rowy := tabtoph + lineh + float64(row)*(rowh+params.RowSpacing)
			c.DrawRectangle(rowx, rowy, colw, rowh)
			c.Fill()
			c.SetColor(color.White)
			c.DrawString(pl.Name, rowx+rowh, rowy+rowh-(params.RowAdditionalHeight)-2)
			var pings string
			switch params.LatencyStyle {
			case LatencyNumber:
				pings = fmt.Sprintf("%d", pl.Ping)
			case LatencyNumberMs:
				pings = fmt.Sprintf("%dms", pl.Ping)
			}
			pw, _ := c.MeasureString(pings)
			c.SetColor(DefaultLatencyColoring(pl.Ping))
			c.DrawString(pings, rowx+colw-pw, rowy+rowh-(params.RowAdditionalHeight)-2)
			if pl.HeadTexture != nil {
				c.DrawImage(imaging.Resize(pl.HeadTexture, int(rowh), int(rowh), imaging.NearestNeighbor), int(rowx), int(rowy))
			}
			plc++
		}
	}
	bottomf := fragmentMessage(c, gg.AlignCenter, *tabbottom, tabw/2, tabtoph+lineh+(rowh+params.RowSpacing)*(float64(params.MaxRows))+lineh*2, params.ChatColorCodes)
	for _, f := range bottomf {
		c.SetColor(f.color)
		c.DrawString(f.str, f.x, f.y)
		if params.DebugTopBottom {
			w, h := c.MeasureString(f.str)
			c.SetColor(color.RGBA{255, 0, 0, 255})
			c.DrawRectangle(f.x, f.y-h, w, h)
			c.Stroke()
		}
	}

	return c.Image()
}
