package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
)
func main() {
	var n int
	var err error
	if len(os.Args) == 1 {
		n = 1
	} else {
		n, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("must provide an number of the amount of guids")
		}
	}

	for i := 0; i < n; i++ {
		fmt.Printf("%s\n",getGuid())
	}

}

func getGuid() string {
	g, err := uuid.NewRandom()
	if err != nil {
	log.Fatalf("unable to generate guid, error: %s", err.Error())
	}
	return g.String()
}