package engine

import (
	"context"

	"github.com/Akmyrat17/carm/models"
	"github.com/Akmyrat17/carm/store"
	"go.opentelemetry.io/otel"
)

type EngineService struct {
	store store.EngineStoreInterface
}

func NewEngineService(store store.EngineStoreInterface) *EngineService {
	return &EngineService{store: store}
}

func (e EngineService) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error) {
	tracer := otel.Tracer("EngineService")
	ctx, span := tracer.Start(ctx, "CreateEngine-Service")
	defer span.End()
	if err := models.ValidateEngineRequest(*engineReq); err != nil {
		return models.Engine{}, err
	}
	engine, err := e.store.CreatedEngine(ctx, engineReq)
	if err != nil {
		return models.Engine{}, err
	}
	return engine, err
}

func (e EngineService) GetEngineById(ctx context.Context, id string) (models.Engine, error) {
	tracer := otel.Tracer("EngineService")
	ctx, span := tracer.Start(ctx, "GeetEngineById-Service")
	defer span.End()
	engine, err := e.store.GetEngineById(ctx, id)
	if err != nil {
		return models.Engine{}, err
	}
	return engine, err
}

func (e EngineService) UpdateEngine(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error) {
	tracer := otel.Tracer("EngineService")
	ctx, span := tracer.Start(ctx, "UpdateEngine-Service")
	defer span.End()
	if err := models.ValidateEngineRequest(*engineReq); err != nil {
		return models.Engine{}, err
	}
	engine, err := e.store.UpdateEngine(ctx, id, engineReq)
	if err != nil {
		return models.Engine{}, err
	}
	return engine, err
}

func (e EngineService) DeleteEngine(ctx context.Context, id string) (models.Engine, error) {
	tracer := otel.Tracer("EngineService")
	ctx, span := tracer.Start(ctx, "DeleteEngine-Service")
	defer span.End()
	engine, err := e.store.DeleteEngine(ctx, id)
	if err != nil {
		return models.Engine{}, err
	}
	return engine, err
}
