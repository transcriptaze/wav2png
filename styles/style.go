package styles

import (
	"encoding/json"
	"image/color"
	"os"

	"github.com/transcriptaze/wav2png/fills"
	"github.com/transcriptaze/wav2png/grids"
	"github.com/transcriptaze/wav2png/kernels"
	"github.com/transcriptaze/wav2png/palettes"
	"github.com/transcriptaze/wav2png/renderers"
	"github.com/transcriptaze/wav2png/renderers/lines"
)

var BLACK = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
var GREEN = color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xff}

type Style struct {
	name     string
	width    uint
	height   uint
	padding  int
	scale    Scale
	fill     Fill
	grid     Grid
	renderer any
}

func NewStyle() Style {
	return Style{
		name:    "default",
		width:   800,
		height:  600,
		padding: 2,

		scale: Scale{
			Horizontal: 1.0,
			Vertical:   1.0,
		},

		fill: Fill{
			Fill:   "solid",
			Colour: "#000000",
			Alpha:  255,
		},

		grid: Grid{
			Grid:   "square",
			Colour: "#008000",
			Alpha:  255,
			Size:   "~64",
			WH:     "~64x48",
		},

		renderer: linesRenderer{
			palette: palette{
				palette: palettes.Ice,
			},
			antialias: kernel{
				kernel: kernels.Vertical,
			},
		},
	}
}

func (s Style) WithWidth(width uint) Style {
	s.width = width

	return s
}

func (s Style) WithHeight(height uint) Style {
	s.height = height

	return s
}

func (s Style) WithPadding(padding int) Style {
	s.padding = padding

	return s
}

func (s Style) WithScale(scale Scale) Style {
	s.scale = scale

	return s
}

func (s Style) WithFill(fill Fill) Style {
	s.fill = fill

	return s
}

func (s Style) WithGrid(grid Grid) Style {
	s.grid = grid

	return s
}

func (s Style) Name() string {
	return s.name
}

func (s Style) Width() uint {
	return s.width
}

func (s Style) Height() uint {
	return s.height
}

func (s Style) Padding() int {
	return s.padding
}

func (s Style) Scale() Scale {
	return s.scale
}

func (s Style) Fill() fills.FillSpec {
	return s.fill.FillSpec()
}

func (s Style) Grid() grids.GridSpec {
	return s.grid.GridSpec()
}

func (s Style) Renderer() renderers.Renderer {
	if l, ok := s.renderer.(*linesRenderer); ok {
		return lines.Lines{
			Palette:   l.palette.Palette(),
			AntiAlias: kernels.Vertical,
		}
	}

	return lines.Lines{
		Palette:   palettes.Ice.Palette(),
		AntiAlias: kernels.Vertical,
	}
}

func (s Style) Load(style string) (Style, error) {
	serializable := struct {
		Name    string         `json:"name"`
		Width   uint           `json:"width"`
		Height  uint           `json:"height"`
		Padding int            `json:"padding"`
		Scale   Scale          `json:"scale"`
		Fill    Fill           `json:"fill"`
		Grid    Grid           `json:"grid"`
		Lines   *linesRenderer `json:"lines"`
	}{
		Width:   s.width,
		Height:  s.height,
		Padding: s.padding,
		Scale:   s.scale,
		Fill:    s.fill,
		Grid:    s.grid,
	}

	if bytes, err := os.ReadFile(style); err != nil {
		return s, err
	} else if err := json.Unmarshal(bytes, &serializable); err != nil {
		return s, err
	} else {
		s.name = serializable.Name
		s.width = serializable.Width
		s.height = serializable.Height
		s.padding = serializable.Padding
		s.scale = serializable.Scale
		s.fill = serializable.Fill
		s.grid = serializable.Grid

		if serializable.Lines != nil {
			s.renderer = serializable.Lines
		}

		return s, nil
	}
}
