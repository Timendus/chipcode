package emulator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/timendus/chipcode/octopus2/src/emulator/silicon8"
)

func Emulate(rom []byte, sequence string) error {
	// We create the emulator here, but we initialize it lazily in the steps
	// below, so we can select the right mode and the thing doesn't complain
	// about the ROM size
	emu := newEmulator()
	emu.rom = rom

	// This implements a little parser for the emulation sequence
	sequence = strings.ReplaceAll(sequence, "\n", ",")
	steps := strings.Split(sequence, ",")
	for _, step := range steps {
		step = strings.ToLower(strings.TrimSpace(step))
		switch {
		case step == "interactive":
			err := emu.init()
			if err != nil {
				return err
			}
			err = emu.interactive()
			if err != nil {
				return err
			}

		case step == "display":
			err := emu.init()
			if err != nil {
				return err
			}
			fmt.Println(emu.display)

		case isNumeric(step):
			err := emu.init()
			if err != nil {
				return err
			}
			cycles, err := strconv.Atoi(step)
			if err != nil {
				return fmt.Errorf("could not parse number in emulation step: '%s'", step)
			}
			frames := cycles / emu.cpf
			for i := 0; i < frames; i++ {
				emu.cpu.ClockTick()
			}
			leftOver := cycles - frames*emu.cpf
			emu.cpu.SetCyclesPerFrame(leftOver)
			emu.cpu.ClockTick()
			emu.cpu.SetCyclesPerFrame(emu.cpf)

		case isStatement(step):
			parts := strings.Split(step, ":")
			Assert(len(parts) >= 2, "Regular expression is not succesfully guarding this case")

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			var params string
			if len(parts) > 2 {
				params = strings.TrimSpace(parts[2])
			}

			// Lazily initialize the emulator unless we're changing settings
			if !(key == "cpf" || key == "mode") {
				err := emu.init()
				if err != nil {
					return err
				}
			}

			switch key {

			// "press: 5"
			case "press":
				k, err := strconv.ParseInt(value, 0, 64)
				if err != nil {
					return fmt.Errorf("could not parse number for key to press in emulation step: '%s'", step)
				}
				if k < 0 || k >= 16 {
					return fmt.Errorf("expected a value from 0 - 15 for key to press in emulation step: '%s'", step)
				}
				emu.cpu.Keyboard[k] = true

			// "release: 0xA"
			case "release":
				k, err := strconv.ParseInt(value, 0, 64)
				if err != nil {
					return fmt.Errorf("could not parse number for key to release in emulation step: '%s'", step)
				}
				if k < 0 || k >= 16 {
					return fmt.Errorf("expected a value from 0 - 15 for key to release in emulation step: '%s'", step)
				}
				emu.cpu.Keyboard[k] = false

			// "save: 0x200: [1 2 3]" (square braces optional)
			case "save":
				k, err := strconv.ParseInt(value, 0, 64)
				if err != nil {
					return fmt.Errorf("could not parse memory address in emulation step: '%s'", step)
				}
				if k < 0 || k >= int64(len(emu.cpu.RAM)) {
					return fmt.Errorf("expected an address inside RAM in emulation step: '%s'", step)
				}
				params = strings.TrimPrefix(params, "[")
				params = strings.TrimSuffix(params, "]")
				data := strings.Fields(params)
				n := len(data)
				if n < 1 || int64(n)+k >= int64(len(emu.cpu.RAM)) {
					return fmt.Errorf("expected data to fit in memory in emulation step: '%s'", step)
				}
				for i, v := range data {
					val, err := strconv.ParseUint(v, 0, 8)
					if err != nil {
						return fmt.Errorf("invalid number '%s' in emulation step: '%s'", v, step)
					}
					emu.cpu.RAM[k+int64(i)] = uint8(val)
				}

			// "load: 0x200: 3"
			case "load":
				k, err := strconv.ParseInt(value, 0, 64)
				if err != nil {
					return fmt.Errorf("could not parse memory address in emulation step: '%s'", step)
				}
				if k < 0 || k >= int64(len(emu.cpu.RAM)) {
					return fmt.Errorf("expected an address inside RAM in emulation step: '%s'", step)
				}
				n, err := strconv.ParseInt(params, 0, 64)
				if err != nil {
					return fmt.Errorf("could not parse number of bytes in emulation step: '%s'", step)
				}
				if n < 1 || n+k >= int64(len(emu.cpu.RAM)) {
					return fmt.Errorf("invalid number of bytes in emulation step: '%s'", step)
				}
				for i, v := range emu.cpu.RAM[k : k+n] {
					fmt.Printf("%04x: %02x\n", k+int64(i), v)
				}

			// "mode: schip"
			case "mode":
				if emu.initialized {
					return fmt.Errorf("can't change the mode on an emulator that's already running in step '%s'", step)
				}
				switch value {
				case "vip":
					emu.mode = silicon8.VIP
				case "blindvip":
					emu.mode = silicon8.BLINDVIP
				case "schip":
					emu.mode = silicon8.SCHIP
				case "xochip":
					emu.mode = silicon8.XOCHIP
				default:
					return fmt.Errorf("invalid emulation mode requested: '%s'. Should be one of 'vip', 'blindvip', 'schip' or 'xochip'", value)
				}

			// "cpf: 150"
			case "cpf":
				cycles, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("could not parse number in cycles per frame setting: '%s'", value)
				}
				if cycles < 1 {
					return fmt.Errorf("cycles per frame should be a positive number, not: '%v'", cycles)
				}
				emu.cpf = cycles
				if emu.initialized {
					emu.cpu.SetCyclesPerFrame(cycles)
				}

			default:
				return fmt.Errorf("unknown statement: '%s' in emulation step '%s'. Should be one of 'press', 'release', 'save', 'load', 'mode' or 'cpf'", key, step)
			}

		default:
			return fmt.Errorf("unknown step in emulation sequence: '%s'", step)
		}
	}

	// Successfully did all the steps!
	return nil
}

func isNumeric(s string) bool {
	re := regexp.MustCompile("^[0-9]+$")
	return re.MatchString(s)
}

func isStatement(s string) bool {
	re := regexp.MustCompile(`^[0-9a-zA-Z]+\s*:\s*[0-9a-zA-Z\-]+(:\s*[0-9a-zA-Z\- \[\]]+)?$`)
	return re.MatchString(s)
}
