package main

import (
	"fmt"
)

func main() {
	go originServer()
	go reverseProxyServer()

	fmt.Println("Servers up and running!")

	for {
		var userInput string
		_, err := fmt.Scanln(&userInput)
		if err != nil {
			return
		}

		if userInput == "quit" || userInput == "exit" || userInput == "q" {
			break
		}
	}
}
