package styles

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed styles
var styles embed.FS

func Load(style string) (any, error) {
	filename := fmt.Sprintf("styles/%v.json", style)
	v := map[string]json.RawMessage{}

	if bytes, err := styles.ReadFile(filename); err != nil {
		return nil, err
	} else if err := json.Unmarshal(bytes, &v); err != nil {
		return nil, err
	} else {
	}

	return nil, nil
}
