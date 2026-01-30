package pkg

import (
	"log"
	"os"
)

func Root() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("root detction failure: %s", err.Error())
	}

	return path
}
