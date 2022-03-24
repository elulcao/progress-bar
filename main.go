package main

import (
	"fmt"
	"time"

	cmd "github.com/elulcao/progress-bar/cmd"
)

func main() {
	pb := cmd.NewPBar()
	pb.SignalHandler()
	pb.Total = 100

	for i := 1; uint16(i) <= pb.Total; i++ {
		pb.RenderPBar(i)
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}

	pb.CleanUp()
}
