package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/freebirdrides/good-beat/beater"
)

func main() {
	err := beat.Run("good-beat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
