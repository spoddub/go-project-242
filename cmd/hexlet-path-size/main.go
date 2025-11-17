package main

import (
	"code"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no path")
	}

	path := os.Args[1]
	size, err := pathsize.GetSize(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d\t%s\n", size, path)
}
