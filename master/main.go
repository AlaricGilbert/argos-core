package main

import (
	api "github.com/AlaricGilbert/argos-core/master/kitex_gen/api/argosmaster"
	"log"
)

func main() {
	svr := api.NewServer(new(ArgosMasterImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
