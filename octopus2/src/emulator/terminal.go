package emulator

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

type terminal struct {
	state *term.State
	input chan byte
}

const (
	KEY_ERROR = iota
	KEY_UP
	KEY_DOWN
	KEY_RIGHT
	KEY_LEFT
	KEY_ESCAPE
)

var relevantKeys = map[int][3]byte{
	KEY_UP:     {27, 91, 65},
	KEY_DOWN:   {27, 91, 66},
	KEY_RIGHT:  {27, 91, 67},
	KEY_LEFT:   {27, 91, 68},
	KEY_ESCAPE: {27, 0, 0},
}

var t terminal

func (t *terminal) init() error {
	fmt.Print("\033[?1049h") // ANSI code to switch to alternate buffer

	var err error
	t.state, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("error setting terminal to raw mode: %s", err)
	}

	t.input = make(chan byte)
	go t.readInput()

	return nil
}

func (t *terminal) deinit() {
	term.Restore(int(os.Stdin.Fd()), t.state)
	fmt.Print("\033[?1049l") // ANSI code to switch back to normal buffer
}

func (t *terminal) readInput() {
	buffer := make([]byte, 3)
outer:
	for {
		num, err := os.Stdin.Read(buffer)
		if err != nil {
			t.input <- KEY_ERROR
		}

		if num == 3 {
			// See if this sequence is in our list of relevant keys with a longer sequence
			for key, bytes := range relevantKeys {
				if buffer[0] == bytes[0] && buffer[1] == bytes[1] && buffer[2] == bytes[2] {
					t.input <- byte(key)
					continue outer
				}
			}

			fmt.Println(buffer)
			time.Sleep(3 * time.Second)
			continue outer
		}

		if num == 1 && (buffer[0] == 27 || buffer[0] == 3) {
			t.input <- KEY_ESCAPE
			continue outer
		}

		// Otherwise just send the first byte
		t.input <- buffer[0]
	}
}

func (t *terminal) clearString() string {
	return "\033[2J\033[H"
}
