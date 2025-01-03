package preprocessor

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
)

func loadImageFile(filename string, modifiers string) (string, error) {
	if _, err := os.Stat(filename); err != nil {
		return "", fmt.Errorf("Requested file '%s' not found", filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("Error reading file '%s': %s", filename, err.Error())
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("Error decoding image '%s': %s", filename, err.Error())
	}
	bounds := img.Bounds()
	if bounds.Size().X%8 != 0 {
		return "", fmt.Errorf("Image width of '%s' is not a multiple of 8 pixels", filename)
	}
	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		rgbaImg = image.NewRGBA(bounds)
		draw.Draw(rgbaImg, bounds, img, bounds.Min, draw.Src)
	}

	label := path.Base(filename)
	label = label[:len(label)-len(path.Ext(filename))]
	mods := parseModifiers(modifiers, bounds.Size().X, bounds.Size().Y)
	imagePlanes := imageToPlanes(rgbaImg, label, mods)
	sprites := make([][]bitmap, len(imagePlanes))
	for i, plane := range imagePlanes {
		sprites[i] = splitIntoSprites(plane, mods.width, mods.height)
	}

	output := ""
	for i := 0; i < len(sprites[0]); i++ {
		for j := 0; j < len(sprites); j++ {
			if mods.labels {
				output += sprites[j][i].label + "\n"
			}
			output += dataToOctoText(sprites[j][i].pixels)
		}
	}

	if mods.debug {
		fmt.Printf("%s - %dx%d\n", path.Base(filename), bounds.Size().X, bounds.Size().Y)
		fmt.Println(mods.string())
		fmt.Println(imagePlanes.string(mods.palette))
		for i := 0; i < len(sprites[0]); i++ {
			colorSprite := planes{}
			for j := 0; j < len(sprites); j++ {
				colorSprite = append(colorSprite, sprites[j][i])
			}
			colorSpriteString := strings.Split(colorSprite.string(mods.palette), "\n")
			planeStrings := make([][]string, 0)
			for j := 0; j < len(sprites); j++ {
				sprite := sprites[j][i]
				fmt.Printf("%s - %dx%d\n", sprite.label, sprite.width, sprite.height)
				planeStrings = append(planeStrings, strings.Split(sprite.string(), "\n"))
			}
			for j := 0; j < len(colorSpriteString); j++ {
				fmt.Print(colorSpriteString[j] + "    ")
				if len(planeStrings) > 1 {
					for k := 0; k < len(planeStrings); k++ {
						fmt.Print(planeStrings[k][j] + "    ")
					}
				}
				fmt.Println()
			}
		}
	}

	return output, nil
}

func parseModifiers(modifiers string, width int, height int) mods {
	// Set sane sprite size defaults
	if width != 16 || height != 16 {
		width = 8
		vertSprites := 1
		for !(height%vertSprites == 0) || height/vertSprites >= 16 {
			vertSprites++
		}
		height = height / vertSprites
	}

	// Is a specific sprite size requested?
	parts := strings.Split(modifiers, " ")
	for _, part := range parts {
		if strings.Index(part, "x") != -1 {
			split := strings.Split(part, "x")
			setWidth, err1 := strconv.Atoi(split[0])
			setHeight, err2 := strconv.Atoi(split[1])
			if err1 != nil || err2 != nil {
				continue
			}
			width = setWidth
			height = setHeight
			break
		}
	}

	// Are we overriding the color palette?
	palette := []color{{0, 0, 0}, {255, 255, 255}}
	if strings.Index(modifiers, "[") != -1 {
		start := strings.Index(modifiers, "[")
		end := strings.Index(modifiers, "]")
		paletteText := modifiers[start+1 : end]
		colors := strings.Split(paletteText, " ")
		palette = []color{}
		for _, color := range colors {
			palette = append(palette, parseColor(color))
		}
	}

	return mods{
		width:    width,
		height:   height,
		palette:  palette,
		labels:   strings.Index(modifiers, "no-labels") == -1,
		debug:    strings.Index(modifiers, "debug") != -1 || strings.Index(modifiers, "debugging") != -1 || strings.Index(modifiers, "verbose") != -1,
		dithered: strings.Index(modifiers, "dithered") != -1 || strings.Index(modifiers, "dither") != -1 || strings.Index(modifiers, "dithering") != -1,
	}
}

func parseColor(input string) color {
	// Remove trailing comma if present
	if strings.Index(input, ",") == len(input)-1 {
		input = input[0 : len(input)-1]
	}

	// Parse as hexadecimal RRBBGG
	if len(input) == 6 {
		red, err1 := strconv.ParseUint(input[0:2], 16, 8)
		green, err2 := strconv.ParseUint(input[2:4], 16, 8)
		blue, err3 := strconv.ParseUint(input[4:6], 16, 8)

		if err1 == nil && err2 == nil && err3 == nil {
			return color{byte(red), byte(green), byte(blue)}
		}
	}

	// Parse as hexadecimal RGB
	if len(input) == 3 {
		red, err1 := strconv.ParseUint(input[0:1]+input[0:1], 16, 8)
		green, err2 := strconv.ParseUint(input[1:2]+input[1:2], 16, 8)
		blue, err3 := strconv.ParseUint(input[2:3]+input[2:3], 16, 8)

		if err1 == nil && err2 == nil && err3 == nil {
			return color{byte(red), byte(green), byte(blue)}
		}
	}

	// Can't parse. Default to black. TODO: better error handling
	return color{0, 0, 0}
}

func imageToPlanes(image *image.RGBA, label string, mods mods) planes {
	if image.Bounds().Size().X%8 != 0 {
		panic("Image width must be a multiple of 8 pixels")
	}
	if mods.palette == nil || len(mods.palette) == 0 {
		panic("Must have a palette")
	}

	if mods.dithered {
		dither(image, mods.palette)
	} else {
		threshold(image, mods.palette)
	}

	width := image.Bounds().Size().X
	height := image.Bounds().Size().Y

	// Create the required planes
	planes := planes{}
	for i := 0; i < int(math.Log2(float64(len(mods.palette)))); i++ {
		planes = append(planes, bitmap{
			label:  label + "-" + strconv.Itoa(i),
			pixels: make([]byte, width*height/8),
			width:  width,
			height: height,
		})
	}
	if len(planes) == 1 {
		// If we only have one plane, labels don't need to include the plane number
		planes[0].label = label
	}

	// Split the colors of the image into single bit planes
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := (y*width + x) * 4
			color := color{
				image.Pix[index+0],
				image.Pix[index+1],
				image.Pix[index+2],
			}
			colorIndex := 0
			for i, col := range mods.palette {
				if color.equals(col) {
					colorIndex = i
				}
			}
			for i, plane := range planes {
				if (colorIndex>>i)&1 == 1 {
					plane.pixels[(y*width+x)/8] |= 1 << (7 - byte(x%8))
				}
			}
		}
	}
	return planes
}

func splitIntoSprites(image bitmap, spriteWidth int, spriteHeight int) []bitmap {
	sprites := []bitmap{}
	for y := 0; y < image.height; y += spriteHeight {
		for x := 0; x < image.width/8; x += spriteWidth / 8 {
			index := y*image.width/8 + x
			pixels := []byte{}
			for rows := 0; rows < spriteHeight; rows++ {
				for cols := 0; cols < spriteWidth/8; cols++ {
					if index+rows*image.width/8+cols < len(image.pixels) {
						pixels = append(pixels, image.pixels[index+rows*image.width/8+cols])
					} else {
						pixels = append(pixels, 0)
					}
				}
			}
			sprites = append(sprites, bitmap{
				label:  fmt.Sprintf(": %s-%d-%d", image.label, x/(spriteWidth/8), y/spriteHeight),
				pixels: pixels,
				width:  spriteWidth,
				height: spriteHeight,
			})
		}
	}
	return sprites
}
