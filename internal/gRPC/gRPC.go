package gRPC

import (
	"Brands/internal/service"
	"Brands/pkg/zerohook"
	"context"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"

	"Brands/internal/api/v1"
)

// Реализация gRPC сервиса
type gRPCServiceServer struct {
	v1.UnimplementedGRPCServiceServer // Встраиваем UnimplementedGRPCServiceServer
	brandService                      *service.BrandService
	log                               zerolog.Logger // Используем логгер из zerohook
	grpcServer                        *grpc.Server   // Экспортируемый grpc.Server
}

// Создание нового gRPC сервиса
func NewgRPCServiceServer(brandService *service.BrandService) *gRPCServiceServer {
	// Инициализация grpc.Server
	grpcServer := grpc.NewServer()
	return &gRPCServiceServer{
		brandService: brandService,
		log:          zerohook.Logger, // Используем глобальный логгер
		grpcServer:   grpcServer,      // Сохраняем grpc.Server
	}
}

// Получить бренд по ID
func (s *gRPCServiceServer) GetBrand(ctx context.Context, req *v1.GetBrandRequest) (*v1.GetBrandResponse, error) {
	brand, err := s.brandService.GetByID(ctx, req.GetBrandId())
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get brand")
		return nil, err
	}

	return &v1.GetBrandResponse{
		BrandId:       brand.ID,
		Name:          brand.Name,
		Description:   brand.Description,
		LogoUrl:       brand.LogoURL,
		CoverImageUrl: brand.CoverImageURL,
		FoundedYear:   int32(brand.FoundedYear),
		OriginCountry: brand.OriginCountry,
		Popularity:    int32(brand.Popularity),
		IsPremium:     brand.IsPremium,
		IsUpcoming:    brand.IsUpcoming,
		CreatedAt:     brand.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Запуск gRPC сервера
func (s *gRPCServiceServer) StartServer(grpcAddr string) error {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		zerohook.Logger.Fatal().Err(err).Msg("Ошибка при запуске gRPC сервера")
		return err
	}

	// Регистрируем сервис и reflection
	v1.RegisterGRPCServiceServer(s.grpcServer, s)
	reflection.Register(s.grpcServer)

	// Логируем запуск gRPC сервера
	zerohook.Logger.Info().Msg("gRPC сервер слушает на " + grpcAddr)

	// Запуск сервера
	if err := s.grpcServer.Serve(lis); err != nil {
		zerohook.Logger.Fatal().Err(err).Msg("Ошибка при запуске gRPC сервера")
		return err
	}

	return nil
}
