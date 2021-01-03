package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan bool)

	newBoard := Initialise(25, 25)

	for {
		fmt.Println(newBoard)
		fmt.Println()

		go newBoard.Iterate(ch)

		time.Sleep(500 * time.Millisecond)

		<-ch
	}
}
