package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/vapourismo/knx-go/knx"
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/dpt"
	"net/http"
	"time"
)

var log = logrus.New()

func main() {
	log.SetLevel(logrus.DebugLevel)

	router := mux.NewRouter()
	router.HandleFunc("/hello", hello)

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		Handler:      router,
	}

	log.Infof("starting webserver on: %s.", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Error(err)
	}

	// Connect to the gateway.
	client, err := knx.NewGroupTunnel("10.0.0.7:3671", knx.DefaultTunnelConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Close upon exiting. Even if the gateway closes the connection, we still have to clean up.
	defer client.Close()

	// Send 20.5°C to group 1/2/3.
	err = client.Send(knx.GroupEvent{
		Command:     knx.GroupWrite,
		Destination: cemi.NewGroupAddr3(1, 2, 3),
		Data:        dpt.DPT_9001(20.5).Pack(),
	})
	if err != nil {
		log.Fatal(err)
	}

}

func hello(w http.ResponseWriter, req *http.Request) {

	log.Debug("entered hello handler")
	fmt.Fprintf(w, "Hello\n")

}
