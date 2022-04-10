package main

import (
	master "github.com/AlaricGilbert/argos-core/master/kitex_gen/master/argosmaster"
	"log"
)

func main() {
	svr := master.NewServer(new(ArgosMasterImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
