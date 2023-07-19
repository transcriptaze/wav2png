package options

import (
	"encoding/json"
	"os"

	"github.com/transcriptaze/wav2png/go/cmd/options"
	"github.com/transcriptaze/wav2png/go/palettes"
	"github.com/transcriptaze/wav2png/go/styles"
)

type settings struct {
	Size      options.Size      `json:"size,omitempty"`
	Palette   palettes.Palette  `json:"palette,omitempty"`
	Fill      styles.Fill       `json:"fill,omitempty"`
	Padding   options.Padding   `json:"padding,omitempty"`
	Grid      styles.Grid       `json:"grid,omitempty"`
	Antialias options.Antialias `json:"antialias,omitempty"`
	Scale     styles.Scale      `json:"scale,omitempty"`
}

func (s *settings) Load(file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, s)
	if err != nil {
		return err
	}

	return nil
}
