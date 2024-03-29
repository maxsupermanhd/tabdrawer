package tabdrawer

import (
	"github.com/maxsupermanhd/go-vmc/v762/chat"
)

// {"text":"\n  §3reddit.com/r/constantiam - patreon.com/constantiam  \n§3phantom@constantiam.net\n"}
// {"text":"\n  §3Ping: 0   TPS: 20   §3Players: 54  \n§6/help - /joindate - /donate\n§6April 2024: $5/$340\n"}

var (
	legacyColorCodes = map[rune]string{
		'0': "black",
		'1': "dark_blue",
		'2': "dark_green",
		'3': "dark_aqua",
		'4': "dark_red",
		'5': "dark_purple",
		'6': "gold",
		'7': "gray",
		'8': "dark_gray",
		'9': "blue",
		'a': "green",
		'b': "aqua",
		'c': "red",
		'd': "light_purple",
		'e': "yellow",
		'f': "white",
	}
)

func cloneMessageStyle(from chat.Message) chat.Message {
	return chat.Message{
		Bold:          from.Bold,
		Italic:        from.Italic,
		UnderLined:    from.UnderLined,
		StrikeThrough: from.StrikeThrough,
		Obfuscated:    from.Obfuscated,
		Color:         from.Color,
	}
}

func ConvertColorCodes(inputString string) chat.Message {
	input := []rune(inputString)
	ret := chat.Message{
		Extra: []chat.Message{},
	}
	building := chat.Message{}
	for i := 0; i < len(input); i++ {
		if input[i] != '§' {
			building.Text += string(input[i])
			continue
		}
		if i+1 >= len(input) {
			continue
		}
		c := input[i+1]
		if c == '§' {
			building.Text += "§"
			i++
			continue
		}
		if c == 'r' {
			ret = ret.Append(building)
			building = chat.Message{}
			i++
			continue
		}
		if c == 'l' {
			ret = ret.Append(building)
			building = cloneMessageStyle(building)
			building.Bold = true
			i++
			continue
		}
		if c == 'o' {
			ret = ret.Append(building)
			building = cloneMessageStyle(building)
			building.Italic = true
			i++
			continue
		}
		if c == 'n' {
			ret = ret.Append(building)
			building = cloneMessageStyle(building)
			building.UnderLined = true
			i++
			continue
		}
		if c == 'm' {
			ret = ret.Append(building)
			building = cloneMessageStyle(building)
			building.StrikeThrough = true
			i++
			continue
		}
		if c == 'k' {
			ret = ret.Append(building)
			building = cloneMessageStyle(building)
			building.Obfuscated = true
			i++
			continue
		}
		if colorcode, ok := legacyColorCodes[c]; ok {
			ret = ret.Append(building)
			building = cloneMessageStyle(building)
			building.Color = colorcode
			i++
			continue
		}
	}
	ret = ret.Append(building)
	return ret
}
