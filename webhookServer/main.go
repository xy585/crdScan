package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

type App struct {
}

func main() {
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}

// StartServer starts the server
func startServer() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "443"
	}

	app := &App{}

	mux := buildRouter(app)

	fmt.Printf("Listening on port %s\n", port)

	return http.ListenAndServeTLS(fmt.Sprintf("0.0.0.0:%s", port), "./cert.pem", "./key.pem", mux)
}

func buildRouter(app *App) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/*", app.handleMutate)
	return r
}

func (app *App) handleMutate(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.Method)
	port := r.URL.Path[1:]
	http.Redirect(w, r, fmt.Sprintf("http://localhost:%s?id=1", port), http.StatusFound) //Construct any HTTP request for apiserver
	fmt.Println(fmt.Sprintf("http://localhost:%s", port))
}
