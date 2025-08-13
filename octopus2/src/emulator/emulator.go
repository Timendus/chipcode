package emulator

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/timendus/chipcode/octopus2/src/emulator/silicon8"
)

type emulator struct {
	cpu         silicon8.CPU
	rom         []byte
	cpf         int
	interacting bool
	display     string
}

func (emu *emulator) init() {
	emu.cpu.RegisterSoundCallbacks(emu.playSound, emu.stopSound)
	emu.cpu.RegisterRandomGenerator(emu.randomByte)
	emu.cpu.RegisterDisplayCallback(emu.render)
	emu.cpu.Reset(silicon8.VIP)
	emu.cpu.SetCyclesPerFrame(emu.cpf)
	emu.cpu.Start()
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

	fmt.Println(t.clearString() + emu.display)
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
	display := ""
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			index := row*width*3 + col*3
			r := buffer[index+0]
			g := buffer[index+1]
			b := buffer[index+2]
			display += fmt.Sprintf("\033[48;2;%v;%v;%vm  ", r, g, b)
		}
		display += "\033[0m\r\n"
	}

	emu.display = display
	if emu.interacting {
		fmt.Println(t.clearString() + display)
	}
}
