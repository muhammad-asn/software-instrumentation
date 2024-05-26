package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	// Load env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		oscall := <-ch
		log.Warn().Msgf("system call:%+v", oscall)
		cancel()
	}()

	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// start: set up any of your logger configuration here if necessary
	setupLogger()
	// end: set up any of your logger configuration here

	// Log the startup message at Info and Debug levels
	log.Info().Msgf("HTTP server successfully started and listening on port %s", server.Addr)
	log.Debug().Msgf("HTTP server configuration: Addr=%s, Handler=%T", server.Addr, server.Handler)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to listen and serve http server")
		}
	}()
	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server gracefully")
	}
}

func setupLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	logFile, _ := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	multi := zerolog.MultiLevelWriter(os.Stdout, logFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Caller().Logger()
	switch logLevel {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	log := log.With().Str("request_id", uuid.New().String()).Logger()
	ctx := log.WithContext(r.Context())
	ctx = addFuncNameToContext(ctx)
	res, err := greeting(ctx, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(res))
}

func greeting(ctx context.Context, name string) (string, error) {
	ctx = addFuncNameToContext(ctx)
	if len(name) < 5 {
		log.Ctx(ctx).Error().Msgf("Failed to process the name because it is too short: %s", name)
		return fmt.Sprintf("Hello %s! Your name is to short\n", name), nil
	}
	log.Ctx(ctx).Debug().Msgf("Processed name successfully: %s", name)
	showTimeRes := showTime(ctx)
	return fmt.Sprintf("Hi %s \n%s", name, showTimeRes), nil
}

func showTime(ctx context.Context) string {
	ctx = addFuncNameToContext(ctx)
	currentTime := time.Now().Format(time.RFC1123)
	log.Ctx(ctx).Debug().Msgf("Current time is: %s", currentTime)
	return fmt.Sprintf("Current time is: %s\n", currentTime)
}

func addFuncNameToContext(ctx context.Context) context.Context {
	pc, _, _, ok := runtime.Caller(1) // Get the caller of addFuncNameToContext
	if !ok {
		log.Ctx(ctx).Warn().Msg("Failed to get function name")
		return ctx
	}
	funcName := runtime.FuncForPC(pc).Name()
	return log.Ctx(ctx).With().Str("func_name", funcName).Logger().WithContext(ctx)
}
