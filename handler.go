package main

import (
	"fn/elasticsearch"
	"fn/moveit"
	"log"
	"net/http"
	"os"
)

func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	indexClient, err := elasticsearch.NewIndexClient[moveit.IndexRecord]("moveit_events")
	if err != nil {
		log.Fatal(err)
	}
	handler, err := moveit.NewMoveItHandler(indexClient)
	if err != nil {
		log.Fatalf("Could not create MoveIt handler: %e", err)
	}

	log.Printf("listening on %s...", listenAddr)
	http.Handle("/api/ingest-moveit", handler)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
