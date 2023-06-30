package options

import (
	"encoding/json"
	"os"

	"github.com/transcriptaze/wav2png/cmd/options"
	"github.com/transcriptaze/wav2png/styles/palettes"
)

type settings struct {
	Size      options.Size      `json:"size,omitempty"`
	Palette   palettes.Palette  `json:"palette,omitempty"`
	Fill      options.Fill      `json:"fill,omitempty"`
	Padding   options.Padding   `json:"padding,omitempty"`
	Grid      options.Grid      `json:"grid,omitempty"`
	Antialias options.Antialias `json:"antialias,omitempty"`
	Scale     options.Scale     `json:"scale,omitempty"`
	Style     string            `json:"style"`
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
