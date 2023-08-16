package main

import (
	"net/http"

	"subruleset/config"
	"subruleset/interfaces"
	"subruleset/tlog"
)

func main() {
	err := config.Watch("config.yaml")
	if err != nil {
		tlog.Error.Fatalf("Failed to watch config: %v", err)
	}
	config := config.Get()

	if config.Token == "" {
		tlog.Error.Fatal("Token is empty")
	}

	mux := interfaces.SetupRouter(config)
	port := ":8080"
	tlog.Info.Printf("Starting server on %s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		tlog.Error.Fatalf("Server failed: %v", err)
	}
}
