package preprocessor

import (
	"image"
	"math"
)

// Threshold an image to the given palette
func threshold(image *image.RGBA, palette []color) {
	width := image.Bounds().Size().X
	height := image.Bounds().Size().Y
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := (y*width + x) * 4

			newColor := nearestColor([]float64{
				float64(image.Pix[index+0]),
				float64(image.Pix[index+1]),
				float64(image.Pix[index+2]),
			}, palette)

			image.Pix[index+0] = newColor.r
			image.Pix[index+1] = newColor.g
			image.Pix[index+2] = newColor.b
		}
	}
}

// Floyd-Steinberg dither an image to given palette
func dither(image *image.RGBA, palette []color) {
	width := image.Bounds().Size().X
	height := image.Bounds().Size().Y

	// We work on a buffer of floats, so we can better preserve the propagated
	// errors. This gives a better resulting image than just clamping to uint8.
	pixelBuffer := make([]float64, width*height*4)
	for i := 0; i < len(image.Pix); i++ {
		pixelBuffer[i] = float64(image.Pix[i])
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := (y*width + x) * 4

			closestColor := nearestColor([]float64{
				pixelBuffer[index+0],
				pixelBuffer[index+1],
				pixelBuffer[index+2],
			}, palette)

			err := []float64{
				pixelBuffer[index+0] - float64(closestColor.r),
				pixelBuffer[index+1] - float64(closestColor.g),
				pixelBuffer[index+2] - float64(closestColor.b),
			}

			pixelBuffer[index+0] = float64(closestColor.r)
			pixelBuffer[index+1] = float64(closestColor.g)
			pixelBuffer[index+2] = float64(closestColor.b)

			// Propagate the difference between the target color and the closest
			// color in the palette to the pixels around this one, that we have
			// yet to process.

			// X + 1
			if x < width-1 {
				pixelBuffer[index+4] += err[0] * float64(7) / 16
				pixelBuffer[index+5] += err[1] * float64(7) / 16
				pixelBuffer[index+6] += err[2] * float64(7) / 16
			}

			if y == height-1 {
				continue
			}
			index += width * 4

			// X - 1, Y + 1
			if x > 0 {
				pixelBuffer[index-4] += err[0] * float64(3) / 16
				pixelBuffer[index-3] += err[1] * float64(3) / 16
				pixelBuffer[index-2] += err[2] * float64(3) / 16
			}

			// X, Y + 1
			pixelBuffer[index+0] += err[0] * float64(5) / 16
			pixelBuffer[index+1] += err[1] * float64(5) / 16
			pixelBuffer[index+2] += err[2] * float64(5) / 16

			// X + 1, Y + 1
			if x < width-1 {
				pixelBuffer[index+4] += err[0] * float64(1) / 16
				pixelBuffer[index+5] += err[1] * float64(1) / 16
				pixelBuffer[index+6] += err[2] * float64(1) / 16
			}
		}
	}

	// Copy back into the image
	for i := 0; i < len(image.Pix); i++ {
		image.Pix[i] = byte(pixelBuffer[i])
	}
}

// Find the closest color in the palette to the given color. Operates on
// float64s to avoid overflows and so it can take in ranges <0 and >255 for the
// dithering algorithm.
func nearestColor(test []float64, palette []color) color {
	distance := math.MaxFloat64
	bestMatch := color{}
	for _, col := range palette {
		// Algorithm for color distance taken from
		// https://stackoverflow.com/questions/2103368/color-logic-algorithm
		rmean := (test[0] + float64(col.r)) / 2
		r := test[0] - float64(col.r)
		g := test[1] - float64(col.g)
		b := test[2] - float64(col.b)
		dist := math.Sqrt(float64((int64((512+rmean)*r*r) >> 8) + int64(4*g*g) + (int64((767-rmean)*b*b) >> 8)))
		if dist < distance {
			distance = dist
			bestMatch = col
		}
	}
	return bestMatch
}
