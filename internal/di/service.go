package di

import (
		"context"
		"log"
		"net/http"
		"zip_downloader/internal/web"
		"go.uber.org/fx"
)

func StartHTTPServer(lc fx.Lifecycle, user_handler *web.UserHandler) {
	// Регистрируем маршруты
	mux := http.NewServeMux()
	web.RegisterRoutes(mux, user_handler) // новый метод, если нужно

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux, 
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Server started on http://localhost:8080")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			return server.Close()
		},
	})
}