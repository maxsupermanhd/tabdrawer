package main

import (
	"bytes"
	"image/png"
	"log"
	"math/rand"
	"os"
	"runtime/debug"

	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
	"github.com/maxsupermanhd/tabdrawer"
)

func main() {
	players := map[uuid.UUID]tabdrawer.TabPlayer{}
	rand.Seed(69421)
	for i := 0; i < 47; i++ {
		players[uuid.New()] = tabdrawer.TabPlayer{
			Name:        generateRandomStirng(4 + rand.Intn(13)),
			Ping:        rand.Intn(500),
			HeadTexture: nil,
			Gamemode:    "survival",
		}
	}

	ctop := &chat.Message{}
	cbottom := &chat.Message{}
	// taken durectly from constantiam
	must(ctop.UnmarshalJSON([]byte(`{"text":"","extra":[{"text":"\n"},{"text":"https://constantiam.net\n","color":"dark_aqua"},{"text":"reddit.com/r/constantiam\n","color":"dark_aqua"},{"text":"phantom@constantiam.net\n","color":"dark_aqua"},{"text":""}]}`)))
	must(cbottom.UnmarshalJSON([]byte(`{"text":"","extra":[{"text":"\n"},{"text":"  "},{"text":"Ping: 14   TPS: ","color":"dark_aqua"},{"text":"20.0   ","color":"green"},{"text":"Players: 42  \n","color":"dark_aqua"},{"text":"/help - /joindate - /donate\n","color":"gold"},{"text":""}]}`)))

	result := tabdrawer.DrawTab(players, ctop, cbottom, nil)
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
