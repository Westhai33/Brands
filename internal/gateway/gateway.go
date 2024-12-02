package gateway

import (
	"context"

	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"Brands/internal/api/v1"
	"Brands/pkg/zerohook"
)

// Запуск gRPC Gateway
func RunGateway(grpcAddress string) error {
	// Создание подключения к gRPC серверу
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure()) // Подключение к gRPC серверу
	if err != nil {
		zerohook.Logger.Fatal().Err(err).Msg("Failed to connect to gRPC server")
		return err
	}
	defer conn.Close()

	// Создание mux для HTTP
	mux := runtime.NewServeMux()

	// Регистрируем gRPC сервис в Gateway
	err = v1.RegisterGRPCServiceHandler(context.Background(), mux, conn)
	if err != nil {
		zerohook.Logger.Fatal().Err(err).Msg("Failed to register gRPC service in gateway")
		return err
	}

	// Запуск HTTP сервера на порту 7000
	httpServer := &http.Server{
		Addr:    ":7000",
		Handler: mux,
	}

	zerohook.Logger.Info().Msg("Starting gRPC Gateway on port 7000")
	return httpServer.ListenAndServe()
}
