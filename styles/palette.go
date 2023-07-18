package styles

import (
	"encoding/json"
	"fmt"

	"github.com/transcriptaze/wav2png/palettes"
	"github.com/transcriptaze/wav2png/wav2png"
)

type palette struct {
	palette palettes.Palette
}

func (p *palette) UnmarshalJSON(bytes []byte) error {
	var s string

	if err := json.Unmarshal(bytes, &s); err == nil {
		switch s {
		case "ice":
			p.palette = palettes.Ice
		case "fire":
			p.palette = palettes.Fire
		case "aurora":
			p.palette = palettes.Aurora
		case "horizon":
			p.palette = palettes.Horizon
		case "amber":
			p.palette = palettes.Amber
		case "blue":
			p.palette = palettes.Blue
		case "green":
			p.palette = palettes.Green
		case "gold":
			p.palette = palettes.Gold
		default:
			p.palette = palettes.Ice
		}
		return nil
	}

	return fmt.Errorf("invalid palette spec")
}

func (p palette) Palette() wav2png.Palette {
	return p.palette.Palette()
}
