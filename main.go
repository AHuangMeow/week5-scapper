package main

import (
	"log"
	"os"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Argument required")
	}

	arg := os.Args[1]

	if arg == "-c" || arg == "--crawl" {
		InitMongoDB("mongodb://localhost:27017", "userdb")
		cookieValue := GetCookie()
		var wg sync.WaitGroup

		for i := 2024000000; i <= 2024999999; i += 50000 {
			wg.Add(1)
			go func(from int) {
				for id := from; id < from+50000; id++ {
					GetUserByID(id, cookieValue)
				}
				wg.Done()
			}(i)
		}

		wg.Wait()
	} else {
		log.Fatal("Invalid Argument")
	}
}
