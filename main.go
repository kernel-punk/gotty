package main

import (
	"github.com/kernel-punk/gotty/gottylib"
	"log"
)

func main() {

	if err := gottylib.Run(gottylib.RunParameters{

		Cmd:  "htop",
		Args: []string{},
		Ssl:  false,
	}); err != nil {
		log.Fatal(err)
	}

}
