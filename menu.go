package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type UserInput struct {
	Action   int
	Key      string
	Filepath string
}

func RenderMenu(s *FileServer) {
	opt := UserInput{
		Action:   0,
		Key:      "",
		Filepath: "",
	}

	for {
		fmt.Println("What do you want to do?")
		fmt.Println("1. Send file")
		fmt.Println("2. Get file")
		fmt.Print("Choose option: ")
		_, err := fmt.Scanf("%d", &opt.Action)
		if err != nil || (opt.Action != 1 && opt.Action != 2) {
			fmt.Println("Invalid input, please choose 1 or 2.")
			RenderDots()
			continue
		}

		switch opt.Action {
		case 1:
			handleSendFiles(s, &opt)
		case 2:
			fmt.Println("Getting...")
			return
		}
	}
}
func handleSendFiles(s *FileServer, opt *UserInput) {
	fmt.Print("Provide key: ")
	fmt.Scanf("%s", &opt.Key)
	for{

	fmt.Print("Provide file with path: ")
	fmt.Scanf("%s", &opt.Filepath)
	fmt.Println("Sending...")
	if len(opt.Key) > 0 && len(opt.Filepath) > 0 {
		fs, err := os.Open(opt.Filepath)
		if err != nil {
			fmt.Printf("Error: %s", err)
			RenderDots()
			continue
		}
		defer fs.Close()
		var fread io.Reader = fs
		s.Store(opt.Key, fread)

		fmt.Println("File sent successfully!")
		return
	}
	}
}

func RenderDots() {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 150)
		fmt.Print(".")
	}
	fmt.Println("\n")
}
