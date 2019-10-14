package main

import (
	"backend/controller"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	server := mux.NewRouter()
	server.HandleFunc("/api/uploadHash", controller.UploadHash).Methods("POST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:8080"},
		AllowCredentials: true,
	})
	handler := c.Handler(server)
	fmt.Println("listening on server 8000")
	err := http.ListenAndServe(":"+port, handler) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
