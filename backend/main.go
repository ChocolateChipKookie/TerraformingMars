package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"terraforming-mars-backend/internal/api"
	"terraforming-mars-backend/internal/database"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticPath, r.URL.Path)

	absStaticPath, err := filepath.Abs(h.staticPath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !strings.HasPrefix(absPath, absStaticPath+string(filepath.Separator)) && absPath != absStaticPath {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	// Initialize database with file path
	db, err := database.Init("data/terraforming_mars.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Create API handler
	apiHandler := api.NewHandler(db)

	// Setup routes
	router := mux.NewRouter()

	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiHandler.RegisterRoutes(apiRouter)

	// Serve static files from dist directory (production build)
	// Check multiple locations: ./dist (production), ../dist (development from backend/)
	// If dist directory exists, serve static files. Otherwise, allow CORS for development.
	distPaths := []string{"./dist", "../dist"}
	var distPath string
	for _, path := range distPaths {
		if _, err := os.Stat(path); err == nil {
			distPath = path
			break
		}
	}

	if distPath != "" {
		log.Println("Serving static files from", distPath)
		spa := spaHandler{staticPath: distPath, indexPath: "index.html"}
		router.PathPrefix("/").Handler(spa)
	} else {
		log.Println("Static files not found, enabling CORS for development")
	}

	// Setup CORS for development
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"}, // React and Vite dev servers
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Start server
	handler := c.Handler(router)
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}