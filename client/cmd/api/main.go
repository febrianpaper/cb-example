package main

import (
	"fbriansyah/client/internal/arangodb"
	"fbriansyah/client/internal/echo"
	"fbriansyah/client/internal/usecase"
	"fbriansyah/client/pkg/httpbreaker"
	"fbriansyah/client/util"
	"log"
	"net/http"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/sony/gobreaker"
)

type Module struct {
	db     driver.Database
	config *util.Config
}

func main() {
	config, err := util.LoadConfig("./")
	if err != nil {
		log.Fatalf(err.Error())
	}

	module := Module{
		config: &config,
	}

	err = module.ArangoInit()
	if err != nil {
		log.Fatalf(err.Error())
	}

	vendorRepo, err := arangodb.NewVendorRepository(module.db, "vendor")
	if err != nil {
		log.Fatalf(`error creating arango client: %v`, err)
	}

	senderLogRepo, err := arangodb.NewSenderlogRepository(module.db, "sender_log")
	if err != nil {
		log.Fatalf(`error creating arango client: %v`, err)
	}
	vendors := []string{"halosis", "infobib"}

	httpClient := &http.Client{}
	httpBreaker := httpbreaker.New(vendors, httpClient, &httpbreaker.Settings{
		Interval: time.Second * 20,
		Timeout:  time.Second * 15,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.Requests >= 3
		},
	})

	uc, err := usecase.New(
		vendors,
		usecase.WithVendorRepo(vendorRepo),
		usecase.WithSenderLogRepo(senderLogRepo),
		usecase.WithHttpClient(httpBreaker),
	)
	if err != nil {
		log.Fatalf(`error creating usecase: %v`, err)
	}

	server := echo.New(uc, ":8888")
	if err := server.Run(); err != nil {
		log.Fatalf(`error running server: %v`, err)
	}
}
