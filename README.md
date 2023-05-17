Bufferio: prototype keyreader module for cli tools.
Takes input from Stdin, processes it, sends it out.

func Bufferio(output chan string)

use:
goroutine set to Bufferio(your_input_ch)
