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
	fmt.Printf("res: %v\n", res)
}
