package main

import (
	"github.com/kernel-punk/gotty/gottyLib"
	"log"
)

func main() {

	if err := gottyLib.Run(gottyLib.RunParameters{

		Cmd:  "htop",
		Args: []string{},
		Ssl:  false,
	}); err != nil {
		log.Fatal(err)
	}

}
