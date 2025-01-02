package preprocessor

import (
	"fmt"
	"strings"
)

type color struct {
	r, g, b byte
}

func (col color) equals(other color) bool {
	return col.r == other.r && col.g == other.g && col.b == other.b
}

func (col color) string() string {
	return fmt.Sprintf("#%02x%02x%02x", col.r, col.g, col.b)
}

type bitmap struct {
	label  string
	pixels []byte
	width  int
	height int
}

func (image *bitmap) string() string {
	output := ""
	stride := image.width / 8
	for y := 0; y < image.height; y++ {
		for i := y * stride; i < (y+1)*stride; i++ {
			output += strings.Replace(
				strings.Replace(
					fmt.Sprintf("%08b", image.pixels[i]),
					"1", "██", -1),
				"0", "  ", -1)
		}
		output += "\n"
	}
	return output
}

type planes []bitmap

func (planes planes) string(palette []color) string {
	if len(planes) == 0 {
		return ""
	}
	output := ""
	for y := 0; y < planes[0].height; y++ {
		for x := 0; x < planes[0].width; x++ {
			paletteIndex := 0
			paletteBit := 1
			for _, plane := range planes {
				if plane.pixels[(y*plane.width+x)/8]&byte(1<<(7-x%8)) != 0 {
					paletteIndex |= paletteBit
				}
				paletteBit <<= 1
			}
			output += fmt.Sprintf(
				"\x1b[38;2;%d;%d;%dm██",
				palette[paletteIndex].r,
				palette[paletteIndex].g,
				palette[paletteIndex].b,
			)
		}
		output += "\x1b[0m\n"
	}
	return output
}

type mods struct {
	width    int
	height   int
	palette  []color
	labels   bool
	debug    bool
	dithered bool
}

func (mods *mods) string() string {
	paletteStrings := []string{}
	for _, col := range mods.palette {
		paletteStrings = append(paletteStrings, col.string())
	}
	return fmt.Sprintf(
		"Modifiers: %dx%d, labels: %t, debug: %t, dithered: %t, palette: %s",
		mods.width,
		mods.height,
		mods.labels,
		mods.debug,
		mods.dithered,
		strings.Join(paletteStrings, ", "),
	)
}
