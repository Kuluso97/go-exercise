package main

import (
	"encoding/json"
	"ex4.12/lib"
	"fmt"
	"log"
	"os"
	"strconv"
)

const usage = `xkcd get N
xkcd index OUTPUT_FILE
xkcd search INDEX_FILE QUERY`

func usageDie() {
	fmt.Println(usage)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	client := comic.NewClient()

	cmd := os.Args[1]
	switch cmd {
	case "get":
		if len(os.Args) != 3 {
			usageDie()
		}

		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "N (%s) must be an int", os.Args[1])
			usageDie()
		}
		comic, err := comic.GetComic(n)
		if err != nil {
			log.Fatal("Error getting comic", err)
		}
		comicJson, _ := json.Marshal(comic)
		fmt.Println(string(comicJson))

	case "build":
		if len(os.Args) != 3 {
			usageDie()
		}

		n, _ := strconv.Atoi(os.Args[2])
		err := comic.BuildIndex(client, n)
		if err != nil {
			log.Fatal("Error serializing indexes", err)
		}

	case "search":
		if len(os.Args) < 3 {
			usageDie()
		}
		result := comic.Search(client, os.Args[2:]...)
		fmt.Println(result)
	}
}
