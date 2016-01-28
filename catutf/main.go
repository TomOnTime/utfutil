package main

// Functions that autodetect UTF-8 and UTF-16/LE/BE but return UTF-8.
// The goal is to create functions that "just do the right thing"
// no matter what UTF encoding is used.

import (
	"fmt"
	"log"
	"strings"

	"github.com/TomOnTime/utfutil"
)

func main() {
	data, err := utfutil.ReadFile("inputfile.txt", utfutil.HTML5)
	if err != nil {
		log.Fatal(err)
	}
	final := strings.Replace(string(data), "\r\n", "\n", -1)
	fmt.Println(final)
}
