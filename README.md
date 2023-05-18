Bufferio: prototype keyreader module for cli tools. 
Takes input from Stdin, processes it, sends it out.

func GetInput(output chan string)

use:
go GetInput(your_input_ch)
var input string := <- your_input_ch