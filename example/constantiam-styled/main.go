package main

import (
	"bytes"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"runtime/debug"

	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"github.com/maxsupermanhd/go-vmc/v762/chat"
	"github.com/maxsupermanhd/tabdrawer"
)

func main() {
	players := map[uuid.UUID]tabdrawer.TabPlayer{}
	var changeUUID uuid.UUID
	rand.Seed(69421)
	for i := 0; i < 47; i++ {
		u := uuid.New()
		if i == 42 {
			changeUUID = u
		}
		players[u] = tabdrawer.TabPlayer{
			Name:        chat.Message{Text: generateRandomStirng(4 + rand.Intn(13))},
			Ping:        rand.Intn(500),
			HeadTexture: nil,
			Gamemode:    "survival",
		}
	}
	players[uuid.New()] = tabdrawer.TabPlayer{
		Name:        chat.Message{Text: generateRandomStirng(4 + rand.Intn(13)), Color: "green"},
		Ping:        rand.Intn(500),
		HeadTexture: nil,
		Gamemode:    "survival",
	}

	params := tabdrawer.TabParameters{
		LatencyColoring:       tabdrawer.DefaultLatencyColoring,
		LatencyStyle:          tabdrawer.LatencyNumberMs,
		ChatColorCodes:        tabdrawer.DefaultChatColorCodes,
		BackgroundColor:       tabdrawer.DefaultDiscordBackgroundColor,
		PlayerBackgroundColor: tabdrawer.DefaultPlayerBackgroundColor,
		RowSpacing:            1,
		ColumnSpacing:         6,
		MaxRows:               20,
		FontColor:             color.White,
		Font:                  truetype.NewFace(noerr(truetype.Parse(noerr(os.ReadFile("./font.ttf")))), &truetype.Options{Size: 16}),
		LineSpacing:           0,
		DebugTopBottom:        false,
		DebugHeight:           false,
		RowAdditionalHeight:   2.0,
		OverridePlayerName: func(u uuid.UUID) *chat.Message {
			if u == changeUUID {
				return &chat.Message{Text: "This is overriden"}
			}
			return nil
		},
	}

	ctop := &chat.Message{}
	cbottom := &chat.Message{}
	// taken durectly from constantiam
	must(ctop.UnmarshalJSON([]byte(`{"text":"","extra":[{"text":"\n"},{"text":"https://constantiam.net\n","color":"dark_aqua"},{"text":"reddit.com/r/constantiam\n","color":"dark_aqua"},{"text":"phantom@constantiam.net\n","color":"dark_aqua"},{"text":""}]}`)))
	must(cbottom.UnmarshalJSON([]byte(`{"text":"","extra":[{"text":"\n"},{"text":"  "},{"text":"Ping: 14   TPS: ","color":"dark_aqua"},{"text":"20.0   ","color":"green"},{"text":"Players: 42  \n","color":"dark_aqua"},{"text":"/help - /joindate - /donate\n","color":"gold"},{"text":""}]}`)))

	result := tabdrawer.DrawTab(players, ctop, cbottom, &params)
	buf := bytes.NewBufferString("")
	must(png.Encode(buf, result))

	must(os.WriteFile("tab.png", buf.Bytes(), 0664))
	// fmt.Print(buf.String())
}

func must(err error) {
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}
}

func noerr[T any](ret T, err error) T {
	must(err)
	return ret
}

func generateRandomStirng(l int) (ret string) {
	for i := 0; i < l; i++ {
		switch rand.Intn(3) {
		case 0:
			ret += string(rune('0' + rand.Intn(9)))
		case 1:
			ret += string(rune('A' + rand.Intn(25)))
		case 2:
			ret += string(rune('a' + rand.Intn(25)))
		}
	}
	return
}
