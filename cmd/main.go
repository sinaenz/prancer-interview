package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"scratch/internal/controller"
	"scratch/internal/service"
)

func main() {

	// service
	coordinatorService := service.NewCenter()

	// router and middleware
	router := mux.NewRouter()

	// controller
	httpController := controller.HttpController{
		CoordinatorService: coordinatorService,
	}
	router.HandleFunc("/move", httpController.Move())

	// run http server
	httpServer := http.Server{
		Addr:    "127.0.0.1:8000",
		Handler: router,
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Panic("http server failed to start")
	}

}
