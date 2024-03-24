package tabdrawer

import (
	"image/color"
	"strings"

	"github.com/fogleman/gg"
	"github.com/maxsupermanhd/go-vmc/v762/chat"
)

type renderFragment struct {
	str   string
	color color.Color
	x, y  float64
}

func concatChatMessage(msg chat.Message) string {
	ret := msg.Text
	for _, e := range msg.Extra {
		ret += concatChatMessage(e)
	}
	return ret
}

func measureMaxLine(c *gg.Context, msg chat.Message, lineinterval float64) (w, h float64) {
	for _, s := range strings.Split(concatChatMessage(msg), "\n") {
		// if len(strings.ReplaceAll(s, " ", "")) == 0 {
		// 	continue
		// }
		ww, hh := c.MeasureString(s)
		if w < ww {
			w = ww
		}
		h += hh + lineinterval
	}
	return
}

func measureChatLine(c *gg.Context, msg chat.Message) (ret bool, w, h float64) {
	strs := strings.Split(msg.Text, "\n")
	w, h = c.MeasureString(strs[0])
	if len(strs) > 1 {
		return true, w, h
	}
	for _, e := range msg.Extra {
		ret, ww, hh := measureChatLine(c, e)
		w += ww
		if ret {
			return true, w, h
		}
		if hh > h {
			h = hh
		}
	}
	return false, w, h
}

func fragmentMessage(c *gg.Context, align gg.Align, msg chat.Message, x, y float64, colorcodes map[string]color.RGBA) []renderFragment {
	lx := float64(0)
	lh := float64(0)
	return fragmentMultilineMessage(c, align, msg, &x, &y, &lx, &lh, 0, 0, colorcodes)
}

func fragmentMultilineMessage(c *gg.Context, align gg.Align, msg chat.Message, x, y, lx, lh *float64, law, lah float64, colorcodes map[string]color.RGBA) []renderFragment {
	col := color.RGBA{255, 255, 255, 255}
	if msg.Color != "" {
		coll, ok := colorcodes[msg.Color]
		if ok {
			col = coll
		}
	}
	c.SetColor(col)
	strs := strings.Split(msg.Text, "\n")
	fragments := []renderFragment{}
	for line := 0; line < len(strs)-1; line++ {
		w, h := c.MeasureString(strs[line])
		var xx float64
		switch align {
		case gg.AlignCenter:
			xx = *x - (*lx+w)/2 + *lx
		case gg.AlignLeft:
			xx = *x + *lx
		case gg.AlignRight:
			xx = *x - w - *lx
		}
		if *lh < h {
			*lh = h
		}
		fragments = append(fragments, renderFragment{
			str:   strs[line],
			color: col,
			x:     xx,
			y:     *y,
		})
		*y += *lh + 3
		*lx = 0
		*lh = 0
	}
	s := strs[len(strs)-1]
	if s != "" {
		w, h := c.MeasureString(s)
		tw := float64(0)
		for _, extra := range msg.Extra {
			brr, ew, eh := measureChatLine(c, extra)
			tw += ew
			if eh > h {
				h = eh
			}
			if brr {
				break
			}
		}
		if *lh < h {
			*lh = h
		}
		var xx float64
		switch align {
		case gg.AlignCenter:
			xx = *x - ((*lx + w + tw + law) / 2) + *lx
		case gg.AlignLeft:
			xx = *x + *lx
		case gg.AlignRight:
			xx = *x - w - *lx
		}
		fragments = append(fragments, renderFragment{
			str:   s,
			color: col,
			x:     xx,
			y:     *y,
		})
		*lx = *lx + w
	}
	for i := 0; i < len(msg.Extra); i++ {
		ew := float64(0)
		eh := float64(0)
		for j := i + 1; j < len(msg.Extra); j++ {
			brr, nw, nh := measureChatLine(c, msg.Extra[j])
			ew += nw
			if eh > nh {
				eh = nh
			}
			if brr {
				break
			}
		}
		fragments = append(fragments, fragmentMultilineMessage(c, align, msg.Extra[i], x, y, lx, lh, ew, eh, colorcodes)...)
	}
	return fragments
}
