package main

import (
	"encoding/json"
	"fmt"

	"github.com/maxsupermanhd/tabdrawer"
)

func main() {
	sample := `\n  §3Ping: 0   TPS: 20   §3Players: 54  \n§6/help - /joindate - /donate\n§6April 2024: $5/$340\n`
	b, _ := tabdrawer.ConvertColorCodes(sample).MarshalJSON()
	i := map[string]any{}
	json.Unmarshal(b, &i)
	r, _ := json.MarshalIndent(i, "", "	")
	fmt.Println(sample)
	fmt.Println(string(r))
}
