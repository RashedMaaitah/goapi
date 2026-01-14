package main

import (
	"fmt"
	"net/http"

	"github.com/RashedMaaitah/goapi/internal/handlers"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetReportCaller(true)
	var r *chi.Mux = chi.NewRouter()
	handlers.Handler(r)

	fmt.Println("Starting GO API service....")

	fmt.Println(`
  _________    ___   ___  ____
 / ___/ __ \  / _ | / _ \/  _/
/ (_ / /_/ / / __ |/ ___// /  
\___/\____/ /_/ |_/_/  /___/  `)

	err := http.ListenAndServe("localhost:8000", r)

	if err != nil {
		log.Error(err)
	}
}
