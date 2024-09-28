package main

import (
	"fmt"
	"log"

	"github.com/pandatix/go-acm/api"
)

func main() {
	cli := api.NewACMClient()
	res, err := cli.Search(&api.SearchParams{
		Request: `(Abstract:"capture the flag")`,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("TOTAL RESULTS: %d\n", res.Results)
	for _, ref := range res.References {
		fmt.Printf("[%s] %s (%s) @ %s\n", ref.Category, ref.Title, ref.PubDate, ref.Conference)
	}
}
