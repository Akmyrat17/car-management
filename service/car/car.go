package car

import (
	"context"

	"github.com/Akmyrat17/carm/models"
	"github.com/Akmyrat17/carm/store"
	"go.opentelemetry.io/otel"
)

type CarService struct {
	store store.CarStoreInterface
}

func NewCarService(store store.CarStoreInterface) *CarService {
	return &CarService{store: store}
}

func (c CarService) GetCarById(ctx context.Context, id string) (models.Car, error) {
	tracer := otel.Tracer("CarService")
	ctx, span := tracer.Start(ctx, "GetCarByd-Handler")
	defer span.End()
	car, err := c.store.GetCarById(ctx, id)
	if err != nil {
		return models.Car{}, err
	}
	return car, err
}

func (c CarService) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	tracer := otel.Tracer("CarService")
	ctx, span := tracer.Start(ctx, "GetCarByBrand-Service")
	defer span.End()

	cars, err := c.store.GetCarByBrand(ctx, brand, isEngine)
	if err != nil {
		return nil, err
	}
	return cars, err
}

func (c CarService) DeleteCar(ctx context.Context, id string) (models.Car, error) {
	tracer := otel.Tracer("CarService")
	ctx, span := tracer.Start(ctx, "DeleteCar-Service")
	defer span.End()

	car, err := c.store.DeleteCar(ctx, id)
	if err != nil {
		return models.Car{}, err
	}
	return car, err
}

func (c CarService) UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (models.Car, error) {
	tracer := otel.Tracer("CarService")
	ctx, span := tracer.Start(ctx, "UpdateCar-Service")
	defer span.End()

	if err := models.CarValidateRequest(*carReq); err != nil {
		return models.Car{}, err
	}
	car, err := c.store.UpdateCar(ctx, id, carReq)
	if err != nil {
		return models.Car{}, err
	}
	return car, err
}

func (c CarService) CreateCar(ctx context.Context, carReq *models.CarRequest) (models.Car, error) {
	tracer := otel.Tracer("CarService")
	ctx, span := tracer.Start(ctx, "CreateCar-Service")
	defer span.End()

	if err := models.CarValidateRequest(*carReq); err != nil {
		return models.Car{}, err
	}
	car, err := c.store.CreateCar(ctx, carReq)
	if err != nil {
		return models.Car{}, err
	}
	return car, err
}
