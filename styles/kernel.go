package styles

import (
	"encoding/json"
	"fmt"

	"github.com/transcriptaze/wav2png/kernels"
)

type kernel struct {
	kernel kernels.Kernel
}

func (k *kernel) UnmarshalJSON(bytes []byte) error {
	var s string

	if err := json.Unmarshal(bytes, &s); err == nil {
		switch s {
		case "none":
			k.kernel = kernels.None

		case "vertical":
			k.kernel = kernels.Vertical

		case "horizontal":
			k.kernel = kernels.Horizontal

		case "soft":
			k.kernel = kernels.Soft

		default:
			k.kernel = kernels.Vertical
		}

		return nil
	}

	return fmt.Errorf("invalid kernel spec")
}

func (k kernel) Kernel() kernels.Kernel {
	return k.kernel
}
