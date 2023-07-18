package styles

import (
	"encoding/json"
)

type linesRenderer struct {
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
