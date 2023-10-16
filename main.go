package main

import (
	"log"
	"time"

	cmd "github.com/elulcao/progress-bar/cmd"
)

func main() {
	pb := cmd.NewPBar()
	pb.CustomMsg = " "
	pb.SignalHandler()
	pb.Total = uint16(10)

	mockLogMessages := []string{
		"Starting the application...",
		"Printing more messages...",
		"sending a message...",
		"receiving a message...",
		"receiving a message...",
		"receiving a message...",
		"receiving a message...",
		"receiving a message...",
		"receiving a message...",
		"receiving a message...",
		"receiving a message...",
		"receiving a message...",
	}

	for i := 1; uint16(i) <= pb.Total; i++ {
		pb.RenderPBar(i)
		pb.CustomMsg = mockLogMessages[i-1]
		log.Println(mockLogMessages[i-1])
		// log.Println(i)               // Do something here
		time.Sleep(5 * time.Second) // Wait 1 second, for demo purpose
	}

	pb.CleanUp()
}
