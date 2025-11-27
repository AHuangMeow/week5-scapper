package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Argument required")
	}

	arg := os.Args[1]

	switch arg {
	case "-c", "--crawl":
		Crawl()
	case "-r", "--reserve":
		Reserve()
	default:
		log.Fatal("Invalid Argument")
	}
}
