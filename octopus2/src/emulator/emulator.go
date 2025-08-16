package emulator

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand/v2"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/timendus/silicon8/src/silicon8"
	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
)

type emulator struct {
	cpu           silicon8.CPU
	rom           []byte
	mode          int
	cpf           int
	initialized   bool
	interacting   bool
	displayWidth  int
	displayHeight int
	displayBuffer *[]byte
}

func newEmulator() emulator {
	return emulator{
		cpu:         silicon8.CPU{},
		mode:        silicon8.VIP,
		cpf:         30,
		initialized: false,
		interacting: false,
	}
}

func (emu *emulator) init() error {
	if emu.initialized {
		return nil
	}
	emu.cpu.RegisterSoundCallbacks(emu.playSound, emu.stopSound)
	emu.cpu.RegisterRandomGenerator(emu.randomByte)
	emu.cpu.RegisterDisplayCallback(emu.render)
	emu.cpu.Reset(emu.mode)
	emu.cpu.SetCyclesPerFrame(emu.cpf)
	err := emu.loadROM()
	if err != nil {
		return err
	}
	emu.cpu.Start()
	emu.initialized = true
	return nil
}

func (emu *emulator) loadROM() error {
	if len(emu.rom)+0x200 > len(emu.cpu.RAM) {
		return fmt.Errorf("ROM file too large to fit in emulator memory")
	}

	for i := 0; i < len(emu.rom); i++ {
		emu.cpu.RAM[i+0x200] = emu.rom[i]
	}
	return nil
}

var keyMap = map[byte]int{
	// Arrow keys
	KEY_UP:    5,
	KEY_DOWN:  8,
	KEY_RIGHT: 9,
	KEY_LEFT:  7,

	// Virtual hex pad
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 0xC,
	'q': 4,
	'w': 5,
	'e': 6,
	'r': 0xD,
	'a': 7,
	's': 8,
	'd': 9,
	'f': 0xE,
	'z': 0xA,
	'x': 0,
	'c': 0xB,
	'v': 0xF,

	// Other number keys
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'0': 0,

	// Special keys
	' ':  6, // Spacebar
	'\r': 4, // Enter
}

func (emu *emulator) interactive() error {
	err := t.init()
	if err != nil {
		return err
	}
	defer t.deinit()

	fmt.Println(t.clearString() + emu.displayToString())
	emu.interacting = true

	for {
		select {
		case key := <-t.input:
			switch key {
			case KEY_ERROR:
				emu.interacting = false
				return fmt.Errorf("could not read input from the terminal")

			case KEY_ESCAPE:
				emu.interacting = false
				return nil

			default:
				val, found := keyMap[key]
				if found {
					emu.triggerKey(val)
				}
			}
		default:
			// Call ClockTick roughly 60 timer per second
			emu.cpu.ClockTick()
			time.Sleep(16 * time.Millisecond)
		}
	}
}

func (emu *emulator) triggerKey(key int) {
	if key < 0 || key > 15 {
		return
	}
	emu.cpu.Keyboard[key] = true
	go func() {
		time.Sleep(100 * time.Millisecond)
		emu.cpu.Keyboard[key] = false
	}()
}

func (emu *emulator) playSound(playing bool, pattern *[16]uint8, pitch float64) {}
func (emu *emulator) stopSound()                                                {}

func (emu *emulator) randomByte() uint8 {
	return uint8(rand.UintN(256))
}

func (emu *emulator) render(width int, height int, buffer []uint8) {
	emu.displayWidth = width
	emu.displayHeight = height
	emu.displayBuffer = &buffer

	if emu.interacting {
		fmt.Println(t.clearString() + emu.displayToString())
	}
}

func (emu *emulator) displayToString() string {
	display := ""
	for row := 0; row < emu.displayHeight; row++ {
		for col := 0; col < emu.displayWidth; col++ {
			index := (row*emu.displayWidth + col) * 3
			r := (*emu.displayBuffer)[index+0]
			g := (*emu.displayBuffer)[index+1]
			b := (*emu.displayBuffer)[index+2]
			if r == 0 && g == 0 && b == 0 {
				// Use a black background and two spaces for black instead. That
				// way, we can copy-paste or pipe the terminal output, ignore
				// the ansi characters and get a somewhat usable image out of
				// it, at least for monochrome roms.
				display += "\033[48;2;0;0;0m  "
			} else {
				display += fmt.Sprintf("\033[38;2;%v;%v;%vm██", r, g, b)
			}
		}
		display += "\033[0m\r\n"
	}
	return display
}

func (emu *emulator) displayToImage() (image.Image, error) {
	if emu.displayBuffer == nil {
		return nil, fmt.Errorf("nothing on the display")
	}

	img := image.NewRGBA(image.Rect(0, 0, emu.displayWidth, emu.displayHeight))

	Assert(emu.displayWidth*emu.displayHeight*3 <= len(*emu.displayBuffer), "We should have at least enough bytes in the display buffer")
	Assert(img.Bounds().Size().X*img.Bounds().Size().Y*4 == len(img.Pix), "We should have exactly the right number of bytes in the target image")

	for row := 0; row < emu.displayHeight; row++ {
		for col := 0; col < emu.displayWidth; col++ {
			imgIndex := (row*emu.displayWidth + col) * 4
			bufIndex := (row*emu.displayWidth + col) * 3
			img.Pix[imgIndex+0] = (*emu.displayBuffer)[bufIndex+0]
			img.Pix[imgIndex+1] = (*emu.displayBuffer)[bufIndex+1]
			img.Pix[imgIndex+2] = (*emu.displayBuffer)[bufIndex+2]
			img.Pix[imgIndex+3] = 0xFF
		}
	}
	return img, nil
}

func (emu *emulator) saveScreenshot(file string, scale int) error {
	var encoders = map[string]func(io.Writer, image.Image) error{
		".jpg":  func(w io.Writer, img image.Image) error { return jpeg.Encode(w, img, nil) },
		".jpeg": func(w io.Writer, img image.Image) error { return jpeg.Encode(w, img, nil) },
		".png":  png.Encode,
		".bmp":  bmp.Encode,
		".gif":  func(w io.Writer, img image.Image) error { return gif.Encode(w, img, nil) },
	}

	extensions := make([]string, len(encoders))
	i := 0
	for ext := range encoders {
		extensions[i] = ext
		i++
	}

	if !slices.Contains(extensions, path.Ext(file)) {
		return fmt.Errorf("can't store image type '%s'. Must be one of '%s'", path.Ext(file), strings.Join(extensions, "', '"))
	}

	img, err := emu.displayToImage()
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("could not store image: %v", err)
	}
	defer f.Close()

	if err = encoders[path.Ext(file)](f, scaleImage(img, scale)); err != nil {
		return fmt.Errorf("failed to encode image: %v", err)
	}
	return nil
}

func scaleImage(src image.Image, scale int) image.Image {
	rect := image.Rect(0, 0, src.Bounds().Size().X*scale, src.Bounds().Size().Y*scale)
	dst := image.NewRGBA(rect)
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	return dst
}
