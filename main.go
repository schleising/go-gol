package main

import (
	"fmt"
	"time"
)

func main() {
	newBoard := Initialise(25, 25)

	for {
		fmt.Println(newBoard)
		fmt.Println()

		time.Sleep(500 * time.Millisecond)

		newBoard.Iterate()
	}
}
