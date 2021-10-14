package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	go enterToChat()
	go spamToChat()
	time.Sleep(1000 * time.Second)
}

func enterToChat() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}

func spamToChat() {
	for {
		fmt.Println("heeeeeej!")
		time.Sleep(1 * time.Second)
	}
}
