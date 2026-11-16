package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/cristianrisueo/graph-bookstore-api/graph"
	"github.com/cristianrisueo/graph-bookstore-api/graph/resolvers"
	"github.com/cristianrisueo/graph-bookstore-api/store"
)

func main() {
	// Logger estructurado que escribe en stdout.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Configuración desde variables de entorno con valores por defecto.
	port := getenv("PORT", "8080")
	dsn := getenv("DB_DSN", "postgres://bookstore:bookstore@localhost:5432/bookstore?sslmode=disable")

	// Conexión a la base de datos.
	st, err := store.New(dsn)
	if err != nil {
		logger.Error("connect to database", "err", err)
		os.Exit(1)
	}
	defer st.Close()

	// Servidor GraphQL con transportes POST y GET e introspección.
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &resolvers.Resolver{Store: st}}))
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})
	srv.Use(extension.Introspection{})
	srv.SetErrorPresenter(resolvers.NewErrorPresenter(logger))
	srv.SetRecoverFunc(resolvers.NewRecoverFunc(logger))

	// Rutas: endpoint GraphQL, playground y health check.
	mux := http.NewServeMux()
	mux.Handle("/query", srv)
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := st.Ping(r.Context()); err != nil {
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Write([]byte("ok"))
	})

	// Servidor HTTP arrancado en segundo plano.
	httpSrv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		logger.Info("server listening", "addr", "http://localhost:"+port)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server", "err", err)
			os.Exit(1)
		}
	}()

	// Apagado ordenado con SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown", "err", err)
	}
}

// getenv devuelve el valor de key o fallback si está vacío.
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
