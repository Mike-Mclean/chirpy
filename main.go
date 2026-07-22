package main

import (
	"chirpy/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"github.com/google/uuid"
	"time"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits 	atomic.Int32
	db 				*database.Queries
	platform		string
}

type User struct {
	ID			uuid.UUID `json:"id"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
	Email     	string    `json:"email"`
}

type Chirp struct {
	ID			uuid.UUID 		`json:"id"`
	CreatedAt	time.Time 		`json:"created_at"`
	UpdatedAt 	time.Time 		`json:"updated_at"`
	Body		string			`json:"body"`
	UserID		uuid.NullUUID	`json:"user_id"`
}

func main() {

	const filepathRoot = "."

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	pf := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		platform: pf,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/users", apiCfg.handlerNewUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.printMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
	Addr: ":8080",
	Handler: mux,
	}

	server.ListenAndServe()
}
