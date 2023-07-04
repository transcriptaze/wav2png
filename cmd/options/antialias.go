package options

import (
	"regexp"
	"strings"

	"github.com/transcriptaze/wav2png/wav2png"
)

type Antialias struct {
	Type string `json:"type"`
}

func (a Antialias) String() string {
	switch a.Type {
	case "none":
		return "none"

	case "horizontal":
		return "horizontal"

	case "vertical":
		return "vertical"

	case "soft":
		return "soft"
	}

	return "??"
}

func (a *Antialias) Set(s string) error {
	ss := strings.ToLower(s)
	match := regexp.MustCompile("^(none|horizontal|vertical|soft)$").FindStringSubmatch(ss)

	if len(match) > 1 {
		switch match[1] {
		case "none":
			a.Type = "none"

		case "horizontal":
			a.Type = "horizontal"

		case "vertical":
			a.Type = "vertical"

		case "soft":
			a.Type = "soft"
		}
	}

	return nil
}

func (a Antialias) Kernel() wav2png.Kernel {
	switch a.Type {
	case "none":
		return wav2png.None

	case "horizontal":
		return wav2png.Horizontal

	case "vertical":
		return wav2png.Vertical

	case "soft":
		return wav2png.Soft
	}

	return wav2png.Soft
}
