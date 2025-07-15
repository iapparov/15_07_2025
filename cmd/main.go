package main

import (
	"zip_downloader/internal/config"
	"zip_downloader/internal/di"
	"zip_downloader/internal/web"

	"go.uber.org/fx"
)



func main() {

	app := fx.New(
		// Провайдеры зависимостей
		fx.Provide(
			config.MustLoad, // подгружаем конфиг
			web.NewUserHandler,
		),

		// Регистрация HTTP-сервера
		fx.Invoke(di.StartHTTPServer),
		fx.Invoke(func(handler *web.UserHandler) {
        handler.StartWorker()
    }),
	)

	app.Run()
}