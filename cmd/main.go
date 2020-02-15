package main

import (
	"log"
	"net/http"
	"os"

	"mongo-auditapi/pkg/api"
	"mongo-auditapi/pkg/config"
	"mongo-auditapi/pkg/db"

	"github.com/rs/cors"
)

func main() {
	c := config.GetConfiguration()
	auditDao, err := db.InitializeDataAccess(c.AuditDBUrl)
	if err != nil {
		log.Printf("ERROR, failed to initialize data access to the Audit database due to error: %v\n", err)
		os.Exit(1)
	}

	a := api.FieldAuditService{Config: c, DataAccess: &db.MongoDBAuditFetcher{Config: c, Dao: auditDao}}
	r := a.InitializeRoutes()
	log.Printf("Server listening on port %s\n", a.Config.APIServicePort)
	log.Fatal(http.ListenAndServe(":"+a.Config.APIServicePort, cors.Default().Handler(r)))
}
