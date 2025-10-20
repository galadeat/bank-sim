package logger

import (
	"fmt"
	"log"
	"os"
)

func Init(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}

	log.SetOutput(f)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return f
}
