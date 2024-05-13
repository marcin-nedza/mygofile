package main

import (
	"fmt"
	"time"
)

func RenderMenu() {
	var opt1 int

	fmt.Println("What do you want to do?")
	fmt.Println("1. Send file")
	fmt.Println("2. Get file")
	fmt.Print("Choose option : ")
	fmt.Scanf("%d", &opt1)

	switch {
	case opt1 == 1:
		fmt.Println("Sending...")
	case opt1 == 2:
		fmt.Println("Geting...")

	default:
		fmt.Println("Wrong answer")
		RenderDots()
		RenderMenu()
	}

}
func RenderDots() {
	fmt.Print("Choose again")
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 150)
		fmt.Print(".")
	}
	fmt.Println("\n")
}
