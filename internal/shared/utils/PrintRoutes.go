package utils

import (
	"fmt"

	"github.com/go-chi/chi/v5"
)

func PrintRoutes(chi *chi.Mux) {
	routes := chi.Routes()

	fmt.Println("Printing Registered Routes")
	for _, path := range routes {
		fmt.Printf("%s ", path.Pattern)
		fmt.Println()
	}
}
