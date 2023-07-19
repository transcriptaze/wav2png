package styles

import (
	"encoding/json"
)

type linesRenderer struct {
	palette   palette
	antialias kernel
}

type columnsRenderer struct {
	barWidth  uint
	barGap    uint
	palette   palette
	antialias kernel
}

func (l *linesRenderer) UnmarshalJSON(bytes []byte) error {
	serializable := struct {
		Palette   json.RawMessage `json:"palette"`
		Antialias json.RawMessage `json:"antialias"`
	}{}

	if err := json.Unmarshal(bytes, &serializable); err != nil {
		return err
	} else {
		palette := palette{}
		kernel := kernel{}

		if err := json.Unmarshal(serializable.Palette, &palette); err != nil {
			return err
		}

		if err := json.Unmarshal(serializable.Antialias, &kernel); err != nil {
			return err
		}

		l.palette = palette
		l.antialias = kernel
	}

	return nil
}

func (c *columnsRenderer) UnmarshalJSON(bytes []byte) error {
	serializable := struct {
		Bar struct {
			Width uint `json:"width"`
			Gap   uint `json:"gap"`
		} `json:"bar"`
		Palette   json.RawMessage `json:"palette"`
		Antialias json.RawMessage `json:"antialias"`
	}{}

	if err := json.Unmarshal(bytes, &serializable); err != nil {
		return err
	} else {
		palette := palette{}
		kernel := kernel{}

		if err := json.Unmarshal(serializable.Palette, &palette); err != nil {
			return err
		}

		if err := json.Unmarshal(serializable.Antialias, &kernel); err != nil {
			return err
		}

		c.barWidth = serializable.Bar.Width
		c.barGap = serializable.Bar.Gap
		c.palette = palette
		c.antialias = kernel
	}

	return nil
}
