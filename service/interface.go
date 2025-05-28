package service

import (
	"context"

	"github.com/Akmyrat17/carm/models"
)

type CarServiceInterface interface {
	GetCarById(ctx context.Context, id string) (models.Car, error)
	GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error)
	DeleteCar(ctx context.Context, id string) (models.Car, error)
	UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (models.Car, error)
	CreateCar(ctx context.Context, carReq *models.CarRequest) (models.Car, error)
}

type EngineServiceInterface interface {
	CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error)
	GetEngineById(ctx context.Context, id string) (models.Engine, error)
	UpdateEngine(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error)
	DeleteEngine(ctx context.Context, id string) (models.Engine, error)
}
