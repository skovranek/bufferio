package bufferio

import (
	exec "golang.org/x/sys/execabs"
	"bufio"
	"os"
	"bytes"
	"fmt"
)
/*
[x] turn off buffer
[x] do not print
[x] go back to normal afterwards

at this point, dev/tty is not writing to terminal,
and each keystroke is being read as byte(s).
What needs to happen is an instantaneous response
to each key, including the arrows.

stdin --> buffer --> process --> buffer --> output (chan)

Process:
[x] process bytes and escape characters
	[x] enter(CR, \n)
	[x] up down
	[x] backspace
	[x] if CR, enter input
	[x] if UP, don't intr, input line to mem, cycle up to previous line
	[ ] Left/Right arrows
*/

func GetInput(output chan string) {
	// turn off buffer
	exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
	// do not print
	exec.Command("stty", "-f", "/dev/tty", "-echo").Run()
	// go back to normal afterwards
	defer func(){exec.Command("stty", "-f", "/dev/tty", "sane").Run()}()

	// stdin --> buffer --> process --> buffer --> output (chan)

	reader := bufio.NewScanner(os.Stdin)
	reader.Split(bufio.ScanBytes)

	history := make([][]byte, 0)
	index := 0

	buffer := make([]byte, 0)
	cursor := 0
	
	logN(7)
	
	// for each byte
	for reader.Scan() {
		// get byte
		b := reader.Bytes()

		// add byte
		buffer = append(buffer, b...)
		
		//log(b)
		cursor++

		//splitter:

		// if byte is NL/CR/enter/'\n'/uint8(10)
			// remove NL byte
			// store buffer to history
			// set index to length
			// send buffer out
			// reset buffer to zero
			// reset cursor
		if i := bytes.IndexByte(buffer, uint8(10)); i >= 0 {
			buffer = removeLast(buffer, 1)
			// if last saved byte slice in history is empty, replace it
			if len(history) > 0 && len(history[len(history)-1]) == 0 {
				history[len(history)-1] = buffer
			} else {
				history = append(history, buffer)
			}
			index = len(history)
			copy := make([]byte, 0)
			copy = append(copy, buffer...)
			output <- string(copy)
			buffer = make([]byte, 0)
			cursor = 0

		// else if byte is BS/DEL/uint8(127)
			// remove byte BS
			// decrement cursor
			// print (without add bytes) left arrow
			// print space (without add bytes)
			// print (without add bytes) left arrow
		} else if i := bytes.IndexByte(buffer, uint8(127)); i >= 0 {
			buffer = removeLast(buffer, 2)
			if cursor > 1 {
				cursor -= 2
				backSpace(1)
			}

		// else if bytes are UP/[]byte{uint8(27), uint8(91), uint8(65)}
			// remove UP bytes
			// remove 3 from cursor
			// counter printed sequence by printing DOWN
			// remove printed line
			// if index == length of history:
				// store buffer to history and decrement index
			// decrement index
			// set buffer to prev string (history[i])
			// set cursor to length of buffer/string
			// print buffer
		} else if i := bytes.Index(buffer, []byte{uint8(27), uint8(91), uint8(65)}); i >= 0 {
			buffer = removeLast(buffer, 3)
			cursor = len(buffer)
			logN(27, 91, 66)
			if index == len(history) && index > 0 {
				history = append(history, buffer)
			}
			if index > 0 {
				backSpace(len(buffer))
				index--
				// copy from history
				buffer = make([]byte, 0)
				buffer = append(buffer, history[index]...)
				cursor = len(buffer)
				log(buffer)
			} else {
				logN(7) // bell
			}

		// else if bytes are DOWN/[]byte{uint8(27), uint8(91), uint8(66)}
			// remove DOWN bytes
			// if index < length of history
				// remove printed sequence
				// increment index
				// set buffer to next string (history[index])
				// set cursor to length of buffer/string
				// print buffer
		} else if i := bytes.Index(buffer, []byte{uint8(27), uint8(91), uint8(66)}); i >= 0 {
			cursor -= 3
			buffer = removeLast(buffer, 3)
			if index < len(history)-1 {
				backSpace(len(buffer))
				index++
				buffer = make([]byte, 0)
				buffer = append(buffer, history[index]...)
				cursor = len(buffer)
				log(buffer)
			} else {
				logN(7) // bell
			}
		}
		// process LEFT/RIGHT arrows?
		fmt.Println()
		log(buffer)
	}
}

func log(data []byte) {
	fmt.Print(string(data))
}

func logN(n ...int) {
	for _, val := range n {
		fmt.Print(string([]byte{uint8(val)}))
	}
}

func removeLast(data []byte, n int) []byte {
	if len(data) >= n {
		return data[0 : len(data)-n]
	}
	return data
}

func backSpace(n int) {
	for i := 0; i < n; i++ {
		logN(27, 91, 68, 32, 27, 91, 68) // LEFT, space, LEFT
	}
}

/*
up:     [27 91 65] string([]byte{uint8(27), uint8(91), uint8(65)})
down:   [27 91 66] string([]byte{uint8(27), uint8(91), uint8(66)})
right:  [27 91 67] string([]byte{uint8(27), uint8(91), uint8(67)})
left:   [27 91 68] string([]byte{uint8(27), uint8(91), uint8(68)})
escape: 27
[:      91
up:     65 (A)
down:   66 (B)
right:  67 (C)
left:   68 (D)
enter:  10
space:  32
tab:     9
*/