package styles

import (
	"fmt"
	"regexp"
	"strconv"
)

type Scale struct {
	Horizontal float64 `json:"horizontal"`
	Vertical   float64 `json:"vertical"`
}

func (s Scale) String() string {
	return fmt.Sprintf("%.1f", s.Vertical)
}

func (s *Scale) Set(v string) error {
	match := regexp.MustCompile(`^([0-9]+(?:\.[0-9]*))$`).FindStringSubmatch(v)

	if len(match) > 1 {
		vscale, err := strconv.ParseFloat(match[1], 64)
		if err == nil && vscale >= 0.2 && vscale <= 5.0 {
			s.Vertical = vscale
		}
	}

	return nil
}
