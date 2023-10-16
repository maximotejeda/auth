package main

import (
	"errors"
	"flag"
	"fmt"

	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/maximotejeda/auth/httpes"
	"github.com/maximotejeda/auth/jwtes"
)

func main() {
	// default flags needed for the different parts of the server to run
	addr := flag.String("addr", "localhost:8080", "address where the server will run")
	timeout := flag.Duration("timeout", 10*time.Second, "timeout per request")
	debugFlag := flag.Int("debug", 0, "flag to add vervose output to service 0 - 1") // log will be added at debug level
	// keyDir := flag.String("keys", "keys", "location where the keys will be stored")
	// dns := flag.String("dns", "file:auth.db?mode=rwc", "database connection string")
	flag.Parse()

	// TODO: initiate key creation and pass to server
	programLevel := new(slog.LevelVar) // info by default
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))
	if *debugFlag > 0 {
		programLevel.Set(slog.LevelDebug)
		slog.Debug("seting debug level to 1")
	}
	logErr := slog.NewLogLogger(slog.NewJSONHandler(os.Stderr, nil), slog.LevelError)

	mux := httpes.NewRegexpSolver()
	mux.Add("GET /health$", httpes.LogMiddleWare(healthHandler))

	mux.Add("GET /metrics$", nil)
	jwt, err := jwtes.NewJwtType()
	if err != nil {
		panic(err)
	}

	claims := map[string]interface{}{
		"id": 1,
		"username": "juan",
		"email": "juan@example.com",
		"rol": "user,admin",
		"sub": "pasword reset",
	}

	j, err := jwt.CreateSignToken( claims ,1*time.Second)
	if err != nil {
		panic(err)
	}
	time.Sleep(995*time.Millisecond)

	clm, err := jwt.RefreshExpiredToken(j)
	if err != nil {
		fmt.Println(clm)
		panic(err)
	}
	fmt.Printf("token: %s \nclaims: %v\n", j, clm)
	server := http.Server{
		Handler:     http.TimeoutHandler(mux, *timeout, "timeout"),
		Addr:        *addr,
		ReadTimeout: *timeout,
		ErrorLog:    logErr,
	}

	slog.Info("starting server", "addr", *addr, "timeout", fmt.Sprintf("%v", timeout), "debug", debugFlag)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logErr.Println("servidor cerro de forma inesprada: ", err)
	}

}


func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("ok"))
	return
}
