package preprocessor

import (
	"bytes"
	"image"
	"reflect"
	"testing"
)

func TestImageToPlanesSinglePlane(t *testing.T) {
	image := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{8, 2},
	})
	image.Pix = []byte{
		0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}
	expected := bitmap{
		label:  "test",
		pixels: []byte{0b01010101, 0b01010101},
		width:  8,
		height: 2,
	}
	planes := imageToPlanes(image, "test", mods{width: 8, height: 2, palette: []color{{0, 0, 0}, {255, 255, 255}}})
	if len(planes) != 1 {
		t.Errorf("Expected 1 plane, got %v", len(planes))
	}
	if !reflect.DeepEqual(expected, planes[0]) {
		t.Errorf("Expected %v, got %v", expected, planes[0])
	}
}

func TestImageToPlanesMultiplePlanes(t *testing.T) {
	image := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{8, 2},
	})
	image.Pix = []byte{
		0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x55, 0x55, 0x55, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xAA, 0xAA, 0xAA, 0xFF,
		0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x55, 0x55, 0x55, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xAA, 0xAA, 0xAA, 0xFF,
	}
	expectedPlane1 := bitmap{
		label:  "test-0",
		pixels: []byte{0b01110100, 0b01110100},
		width:  8,
		height: 2,
	}
	expectedPlane2 := bitmap{
		label:  "test-1",
		pixels: []byte{0b01010101, 0b01010101},
		width:  8,
		height: 2,
	}
	planes := imageToPlanes(image, "test", mods{width: 8, height: 2, palette: []color{{0, 0, 0}, {0x55, 0x55, 0x55}, {0xAA, 0xAA, 0xAA}, {255, 255, 255}}})
	if len(planes) != 2 {
		t.Errorf("Expected 2 planes, got %v", len(planes))
	}
	if !reflect.DeepEqual(expectedPlane1, planes[0]) {
		t.Errorf("Expected %v, got %v", expectedPlane1, planes[0])
	}
	if !reflect.DeepEqual(expectedPlane2, planes[1]) {
		t.Errorf("Expected %v, got %v", expectedPlane2, planes[1])
	}
}

func TestParseModifiersDefaultsToSaneDimensions(t *testing.T) {
	mods := parseModifiers("", 8, 8)
	if mods.width != 8 || mods.height != 8 {
		t.Errorf("Expected width and height to be 8, got %v", mods)
	}

	mods = parseModifiers("", 16, 16)
	if mods.width != 16 || mods.height != 16 {
		t.Errorf("Expected width and height to be 16, got %v", mods)
	}

	mods = parseModifiers("", 8, 16)
	if mods.width != 8 || mods.height != 8 {
		t.Errorf("Expected width and height to be 8, got %v", mods)
	}

	mods = parseModifiers("", 16, 8)
	if mods.width != 8 || mods.height != 8 {
		t.Errorf("Expected width and height to be 8, got %v", mods)
	}

	mods = parseModifiers("", 8, 12)
	if mods.width != 8 || mods.height != 12 {
		t.Errorf("Expected width and height to be 8, 12, got %v", mods)
	}

	mods = parseModifiers("", 8, 24)
	if mods.width != 8 || mods.height != 12 {
		t.Errorf("Expected width and height to be 8, 12, got %v", mods)
	}

	mods = parseModifiers("", 8, 36)
	if mods.width != 8 || mods.height != 12 {
		t.Errorf("Expected width and height to be 8, 12, got %v", mods)
	}

	mods = parseModifiers("", 16, 24)
	if mods.width != 8 || mods.height != 12 {
		t.Errorf("Expected width and height to be 8, 12, got %v", mods)
	}

	mods = parseModifiers("", 16, 36)
	if mods.width != 8 || mods.height != 12 {
		t.Errorf("Expected width and height to be 8, 12, got %v", mods)
	}

	mods = parseModifiers("", 64, 6)
	if mods.width != 8 || mods.height != 6 {
		t.Errorf("Expected width and height to be 8, 6, got %v", mods)
	}

	mods = parseModifiers("", 32, 16)
	if mods.width != 8 || mods.height != 8 {
		t.Errorf("Expected width and height to be 8, 8, got %v", mods)
	}
}

func TestParseModifiersAcceptsDimensions(t *testing.T) {
	mods := parseModifiers("8x8", 64, 6)
	if mods.width != 8 || mods.height != 8 {
		t.Errorf("Expected width and height to be 8, 8, got %v", mods)
	}

	mods = parseModifiers("16x16", 32, 16)
	if mods.width != 16 || mods.height != 16 {
		t.Errorf("Expected width and height to be 16, 16, got %v", mods)
	}

	mods = parseModifiers("8x4", 16, 8)
	if mods.width != 8 || mods.height != 4 {
		t.Errorf("Expected width and height to be 8, 4, got %v", mods)
	}
}

func TestParseModifiersShowsLabels(t *testing.T) {
	mods := parseModifiers("", 8, 8)
	if !mods.labels {
		t.Errorf("Expected labels to be true, got %v", mods)
	}

	mods = parseModifiers("8x8", 8, 8)
	if !mods.labels {
		t.Errorf("Expected labels to be true, got %v", mods)
	}

	mods = parseModifiers("labels", 8, 8)
	if !mods.labels {
		t.Errorf("Expected labels to be true, got %v", mods)
	}

	mods = parseModifiers("no-labels", 8, 8)
	if mods.labels {
		t.Errorf("Expected labels to be false, got %v", mods)
	}
}

func TestParseModifiersUnderstandsColors(t *testing.T) {
	mods := parseModifiers("[414141, 197eb3, e59823, ffffff]", 8, 8)
	if !mods.palette[0].equals(color{0x41, 0x41, 0x41}) {
		t.Errorf("Expected first color to be 414141, got %v", mods.palette[0])
	}
	if !mods.palette[1].equals(color{0x19, 0x7e, 0xb3}) {
		t.Errorf("Expected second color to be 197eb3, got %v", mods.palette[1])
	}
	if !mods.palette[2].equals(color{0xe5, 0x98, 0x23}) {
		t.Errorf("Expected third color to be e59823, got %v", mods.palette[2])
	}
	if !mods.palette[3].equals(color{0xff, 0xff, 0xff}) {
		t.Errorf("Expected fourth color to be ffffff, got %v", mods.palette[3])
	}

	mods = parseModifiers("[444 17b e92 fff]", 8, 8)
	if !mods.palette[0].equals(color{0x44, 0x44, 0x44}) {
		t.Errorf("Expected first color to be 444444, got %v", mods.palette[0])
	}
	if !mods.palette[1].equals(color{0x11, 0x77, 0xbb}) {
		t.Errorf("Expected second color to be 1177bb, got %v", mods.palette[1])
	}
	if !mods.palette[2].equals(color{0xee, 0x99, 0x22}) {
		t.Errorf("Expected third color to be ee9922, got %v", mods.palette[2])
	}
	if !mods.palette[3].equals(color{0xff, 0xff, 0xff}) {
		t.Errorf("Expected fourth color to be ffffff, got %v", mods.palette[3])
	}
}

func TestSplitIntoSpritesSingleSprite(t *testing.T) {
	image := bitmap{
		label:  "test",
		width:  8,
		height: 6,
		pixels: []byte{1, 2, 3, 4, 5, 6},
	}

	sprites := splitIntoSprites(image, 8, 6)
	if len(sprites) != 1 {
		t.Errorf("Expected 1 sprite, got %d", len(sprites))
	}

	sprite := sprites[0]
	if sprite.label != ": test-0-0" {
		t.Errorf("Expected label to be ': test-0-0', got '%s'", sprite.label)
	}
	if sprite.width != 8 {
		t.Errorf("Expected width to be 8, got %d", sprite.width)
	}
	if sprite.height != 6 {
		t.Errorf("Expected height to be 6, got %d", sprite.height)
	}
	if !bytes.Equal(sprite.pixels, []byte{1, 2, 3, 4, 5, 6}) {
		t.Errorf("Expected pixels to be [1, 2, 3, 4, 5, 6], got %v", sprite.pixels)
	}
}

func TestSplitIntoSpritesVerticalSplit(t *testing.T) {
	image := bitmap{
		label:  "test",
		width:  8,
		height: 6,
		pixels: []byte{1, 2, 3, 4, 5, 6},
	}

	sprites := splitIntoSprites(image, 8, 3)
	if len(sprites) != 2 {
		t.Errorf("Expected 2 sprites, got %d", len(sprites))
	}

	sprite := sprites[0]
	if sprite.label != ": test-0-0" {
		t.Errorf("Expected label to be ': test-0-0', got '%s'", sprite.label)
	}
	if sprite.width != 8 {
		t.Errorf("Expected width to be 8, got %d", sprite.width)
	}
	if sprite.height != 3 {
		t.Errorf("Expected height to be 3, got %d", sprite.height)
	}
	if !bytes.Equal(sprite.pixels, []byte{1, 2, 3}) {
		t.Errorf("Expected pixels to be [1, 2, 3], got %v", sprite.pixels)
	}

	sprite = sprites[1]
	if sprite.label != ": test-0-1" {
		t.Errorf("Expected label to be ': test-0-1', got '%s'", sprite.label)
	}
	if sprite.width != 8 {
		t.Errorf("Expected width to be 8, got %d", sprite.width)
	}
	if sprite.height != 3 {
		t.Errorf("Expected height to be 3, got %d", sprite.height)
	}
	if !bytes.Equal(sprite.pixels, []byte{4, 5, 6}) {
		t.Errorf("Expected pixels to be [4, 5, 6], got %v", sprite.pixels)
	}
}

func TestSplitIntoSpritesHorizontalSplit(t *testing.T) {
	image := bitmap{
		label:  "test",
		width:  16,
		height: 3,
		pixels: []byte{1, 2, 3, 4, 5, 6},
	}

	sprites := splitIntoSprites(image, 8, 3)
	if len(sprites) != 2 {
		t.Errorf("Expected 2 sprites, got %d", len(sprites))
	}

	sprite := sprites[0]
	if sprite.label != ": test-0-0" {
		t.Errorf("Expected label to be ': test-0-0', got '%s'", sprite.label)
	}
	if sprite.width != 8 {
		t.Errorf("Expected width to be 8, got %d", sprite.width)
	}
	if sprite.height != 3 {
		t.Errorf("Expected height to be 3, got %d", sprite.height)
	}
	if !bytes.Equal(sprite.pixels, []byte{1, 3, 5}) {
		t.Errorf("Expected pixels to be [1, 3, 5], got %v", sprite.pixels)
	}

	sprite = sprites[1]
	if sprite.label != ": test-1-0" {
		t.Errorf("Expected label to be ': test-1-0', got '%s'", sprite.label)
	}
	if sprite.width != 8 {
		t.Errorf("Expected width to be 8, got %d", sprite.width)
	}
	if sprite.height != 3 {
		t.Errorf("Expected height to be 3, got %d", sprite.height)
	}
	if !bytes.Equal(sprite.pixels, []byte{2, 4, 6}) {
		t.Errorf("Expected pixels to be [2, 4, 6], got %v", sprite.pixels)
	}
}

func TestSplitIntoSpriteBothDirections(t *testing.T) {
	image := bitmap{
		label:  "test",
		width:  16,
		height: 6,
		pixels: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
	}

	sprites := splitIntoSprites(image, 8, 3)
	if len(sprites) != 4 {
		t.Errorf("Expected 4 sprites, got %d", len(sprites))
	}

	sprite := sprites[0]
	if sprite.label != ": test-0-0" {
		t.Errorf("Expected label to be ': test-0-0', got '%s'", sprite.label)
	}
	if sprite.width != 8 {
		t.Errorf("Expected width to be 8, got %d", sprite.width)
	}
	if sprite.height != 3 {
		t.Errorf("Expected height to be 3, got %d", sprite.height)
	}
	if !bytes.Equal(sprite.pixels, []byte{1, 3, 5}) {
		t.Errorf("Expected pixels to be [1, 3, 5], got %v", sprite.pixels)
	}

	sprite = sprites[1]
	if sprite.label != ": test-1-0" {
		t.Errorf("Expected label to be ': test-1-0', got '%s'", sprite.label)
	}
	if sprite.width != 8 {
		t.Errorf("Expected width to be 8, got %d", sprite.width)
	}
	if sprite.height != 3 {
		t.Errorf("Expected height to be 3, got %d", sprite.height)
	}
	if !bytes.Equal(sprite.pixels, []byte{2, 4, 6}) {
		t.Errorf("Expected pixels to be [2, 4, 6], got %v", sprite.pixels)
	}

	sprite = sprites[2]
	if sprite.label != ": test-0-1" {
		t.Errorf("Expected label to be ': test-0-1', got '%s'", sprite.label)
	}
	if sprite.width != 8 {
		t.Errorf("Expected width to be 8, got %d", sprite.width)
	}
	if sprite.height != 3 {
		t.Errorf("Expected height to be 3, got %d", sprite.height)
	}
	if !bytes.Equal(sprite.pixels, []byte{7, 9, 11}) {
		t.Errorf("Expected pixels to be [7, 9, 11], got %v", sprite.pixels)
	}

	sprite = sprites[3]
	if sprite.label != ": test-1-1" {
		t.Errorf("Expected label to be ': test-1-1', got '%s'", sprite.label)
	}
	if sprite.width != 8 {
		t.Errorf("Expected width to be 8, got %d", sprite.width)
	}
	if sprite.height != 3 {
		t.Errorf("Expected height to be 3, got %d", sprite.height)
	}
	if !bytes.Equal(sprite.pixels, []byte{8, 10, 12}) {
		t.Errorf("Expected pixels to be [8, 10, 12], got %v", sprite.pixels)
	}
}
