package main

import (
	"log"
	"os"

	"github.com/abusquets/book_management_api_with_gqlgen/pkg/cmd/booking"
)

func main() {
	if err := booking.NewCommand(os.Args); err != nil {
		log.Fatal(err)
	}

}
