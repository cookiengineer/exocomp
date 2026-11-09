package tty

import "fmt"
import "strconv"
import "strings"
import "sync"
import "time"
import "unicode/utf8"

const escapeTimeout = 25 * time.Millisecond

type Input struct {
	terminal *Terminal
	bytes    chan byte
	errors   chan error
	once     *sync.Once
	closed   chan struct{}
}

func NewInput(terminal *Terminal) *Input {

	input := &Input{
		terminal: terminal,
		bytes:    make(chan byte, 4096),
		errors:   make(chan error, 1),
		once:     &sync.Once{},
		closed:   make(chan struct{}),
	}

	input.start()

	return input

}

func (input *Input) Close() {

	input.once.Do(func() {
		close(input.closed)
	})

}

func (input *Input) start() {

	go func() {

		buffer := make([]byte, 1024)

		for {

			count, err := input.terminal.Input.Read(buffer)

			if count > 0 {

				for index := 0; index < count; index++ {

					select {

					case input.bytes <- buffer[index]:
					case <-input.closed:
						return

					}

				}

			}

			if err != nil {

				select {

				case input.errors <- err:
				default:

				}

				return

			}

		}

	}()

}

func (input *Input) Read() (Key, error) {

	first, err := input.next(0)

	if err != nil {
		return Key{Type: KeyEOF}, err
	}

	if first == 0x1b {
		return input.readEscape()
	}

	return input.readSimple(first)

}

func (input *Input) readSimple(value byte) (Key, error) {

	if value == 0x0d {
		return Key{Type: KeyEnter}, nil
	}

	if value == 0x09 {
		return Key{Type: KeyTab}, nil
	}

	if value == 0x7f || value == 0x08 {
		return Key{Type: KeyBackspace}, nil
	}

	if value == 0x03 {
		return Key{Type: KeyRune, Rune: 'c', Ctrl: true}, nil
	}

	if value == 0x04 {
		return Key{Type: KeyRune, Rune: 'd', Ctrl: true}, nil
	}

	if value >= 0x01 && value <= 0x1a {
		return Key{Type: KeyRune, Rune: rune('a' + value - 1), Ctrl: true}, nil
	}

	if value < 0x20 {
		return Key{Type: KeyUnknown}, nil
	}

	if value < 0x80 {
		return Key{Type: KeyRune, Rune: rune(value)}, nil
	}

	return input.readUTF8(value)

}

func (input *Input) readUTF8(first byte) (Key, error) {

	length := 1

	if first&0xe0 == 0xc0 {
		length = 2
	} else if first&0xf0 == 0xe0 {
		length = 3
	} else if first&0xf8 == 0xf0 {
		length = 4
	}

	buffer := make([]byte, 0, length)
	buffer = append(buffer, first)

	for len(buffer) < length {

		next, err := input.next(escapeTimeout)

		if err != nil {
			return Key{Type: KeyUnknown}, nil
		}

		buffer = append(buffer, next)

	}

	if utf8.Valid(buffer) {

		value, _ := utf8.DecodeRune(buffer)
		return Key{Type: KeyRune, Rune: value}, nil

	} else {
		return Key{Type: KeyUnknown}, nil
	}

}

func (input *Input) readEscape() (Key, error) {

	next, err := input.next(escapeTimeout)

	if err != nil {
		return Key{Type: KeyEscape}, nil
	}

	if next == '[' {
		return input.readCSI()
	}

	if next == 'O' {
		return input.readSS3()
	}

	key, err0 := input.readSimple(next)

	if err0 == nil {

		key.Alt = true
		return key, nil

	} else {
		return Key{Type: KeyEscape}, nil
	}

}

func (input *Input) readCSI() (Key, error) {

	parameters := strings.Builder{}
	intermediate := byte(0)

	for {

		next, err := input.next(escapeTimeout)

		if err != nil {
			return Key{Type: KeyEscape}, nil
		}

		if next >= 0x40 && next <= 0x7e {
			return input.decodeCSI(parameters.String(), intermediate, next)
		}

		if next >= 0x20 && next <= 0x2f {
			intermediate = next
		} else {
			parameters.WriteByte(next)
		}

	}

}

func (input *Input) readSS3() (Key, error) {

	next, err := input.next(escapeTimeout)

	if err != nil {
		return Key{Type: KeyEscape}, nil
	}

	switch next {

	case 'A':
		return Key{Type: KeyArrowUp}, nil
	case 'B':
		return Key{Type: KeyArrowDown}, nil
	case 'C':
		return Key{Type: KeyArrowRight}, nil
	case 'D':
		return Key{Type: KeyArrowLeft}, nil
	case 'H':
		return Key{Type: KeyHome}, nil
	case 'F':
		return Key{Type: KeyEnd}, nil
	default:
		return Key{Type: KeyUnknown}, nil

	}

}

func (input *Input) decodeCSI(parameters string, intermediate byte, final byte) (Key, error) {

	if final == '~' {

		if parameters == "200" || strings.HasPrefix(parameters, "200;") {
			return input.readPaste()
		}

		if strings.HasPrefix(parameters, "27;") {
			return decodeModifyOtherKeys(parameters)
		}

		code, err := strconv.Atoi(strings.Split(parameters, ";")[0])

		if err == nil {

			switch code {

			case 1, 7:
				return Key{Type: KeyHome}, nil
			case 2:
				return Key{Type: KeyInsert}, nil
			case 3:
				return Key{Type: KeyDelete}, nil
			case 4, 8:
				return Key{Type: KeyEnd}, nil
			case 5:
				return Key{Type: KeyPageUp}, nil
			case 6:
				return Key{Type: KeyPageDown}, nil
			default:
				return Key{Type: KeyUnknown}, nil

			}

		} else {
			return Key{Type: KeyUnknown}, nil
		}

	}

	if final == 'u' {
		return decodeCSIU(parameters)
	}

	parts := strings.Split(parameters, ";")
	modifier := 0

	if len(parts) > 1 {

		parsed, err := strconv.Atoi(parts[1])
		if err == nil {
			modifier = parsed
		}

	}

	shift, alt, ctrl := decodeModifier(modifier)

	switch final {

	case 'A':
		return Key{Type: KeyArrowUp, Shift: shift, Alt: alt, Ctrl: ctrl}, nil
	case 'B':
		return Key{Type: KeyArrowDown, Shift: shift, Alt: alt, Ctrl: ctrl}, nil
	case 'C':
		return Key{Type: KeyArrowRight, Shift: shift, Alt: alt, Ctrl: ctrl}, nil
	case 'D':
		return Key{Type: KeyArrowLeft, Shift: shift, Alt: alt, Ctrl: ctrl}, nil
	case 'H':
		return Key{Type: KeyHome, Shift: shift, Alt: alt, Ctrl: ctrl}, nil
	case 'F':
		return Key{Type: KeyEnd, Shift: shift, Alt: alt, Ctrl: ctrl}, nil
	case 'Z':
		return Key{Type: KeyTab, Shift: true}, nil
	default:
		return Key{Type: KeyUnknown}, nil

	}

}

func (input *Input) readPaste() (Key, error) {

	builder := strings.Builder{}
	sequence := make([]byte, 0)

	for {

		next, err := input.next(0)

		if err != nil {
			return Key{Type: KeyPaste, Text: builder.String()}, nil
		}

		sequence = append(sequence, next)

		if len(sequence) >= 6 {

			if string(sequence[len(sequence)-6:]) == "\033[201~" {

				value := builder.String()
				value = strings.TrimSuffix(value, "\033[201~")

				return Key{Type: KeyPaste, Text: value}, nil

			}

		}

		builder.WriteByte(next)

	}

}

func (input *Input) next(timeout time.Duration) (byte, error) {

	if timeout <= 0 {

		select {

		case value := <-input.bytes:
			return value, nil
		case err := <-input.errors:
			return 0, err
		case <-input.closed:
			return 0, fmt.Errorf("tty.Input.next: Input closed.")

		}

	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {

	case value := <-input.bytes:
		return value, nil
	case err := <-input.errors:
		return 0, err
	case <-timer.C:
		return 0, fmt.Errorf("tty.Input.next: Input timeout.")
	case <-input.closed:
		return 0, fmt.Errorf("tty.Input.next: Input closed.")

	}

}
