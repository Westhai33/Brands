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
}

// Создание нового gRPC сервиса
func NewgRPCServiceServer(brandService *service.BrandService) *gRPCServiceServer {
	return &gRPCServiceServer{
		brandService: brandService,
		log:          zerohook.Logger, // Используем глобальный логгер
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
		UpdatedAt:     brand.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// Запуск gRPC сервера
func RunGRPCServer(brandService *service.BrandService, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		zerohook.Logger.Fatal().Err(err).Msg("Failed to listen on port")
		return err
	}

	grpcServer := grpc.NewServer()

	// Регистрируем наш сервис
	v1.RegisterGRPCServiceServer(grpcServer, NewgRPCServiceServer(brandService))

	// Регистрируем рефлексию (опционально, для удобства)
	reflection.Register(grpcServer)

	// Логирование запуска сервера
	zerohook.Logger.Info().Msgf("gRPC server is running on %s", port)

	// Запуск сервера
	return grpcServer.Serve(lis)
}
