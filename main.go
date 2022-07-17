package main

import (
	"fmt"
	"time"

	cmd "github.com/elulcao/progress-bar/cmd"
)

func main() {
	pb := cmd.NewPBar()
	pb.SignalHandler()
	pb.Total = uint16(10)

	for i := 1; uint16(i) <= pb.Total; i++ {
		pb.RenderPBar(i)
		fmt.Println(i)              // Do something here
		time.Sleep(1 * time.Second) // Wait 1 second, for demo purpose
	}

	pb.CleanUp()
}
