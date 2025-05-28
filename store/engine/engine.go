package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Akmyrat17/carm/models"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type EngineStore struct {
	db *sql.DB
}

func New(db *sql.DB) *EngineStore {
	return &EngineStore{db: db}
}

func (e EngineStore) GetEngineById(ctx context.Context, id string) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")
	ctx, span := tracer.Start(ctx, "GetEngineById-Store")
	defer span.End()
	var engine models.Engine
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return engine, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("error rolling back transaction: %v", rbErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				fmt.Printf("error committing transaction: %v", commitErr)
			}
		}
	}()
	err = tx.QueryRowContext(ctx, "SELECT id, displacement, no_of_cylinders, car_range FROM engine WHERE id = $1", id).Scan(&engine.ID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return engine, errors.New("engine not found in database")
		}
		return engine, err
	}
	return engine, nil
}

func (e EngineStore) CreatedEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")
	ctx, span := tracer.Start(ctx, "CerateEngine-Store")
	defer span.End()
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("error rolling back transaction: %v", rbErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				fmt.Printf("error committing transaction: %v", commitErr)
			}
		}
	}()

	engineId := uuid.New()

	createdAt := time.Now()
	updatedAt := createdAt

	_, err = tx.ExecContext(ctx,
		"INSERT INTO engine (id, displacement, no_of_cylinders, car_range, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", engineId, engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange, createdAt, updatedAt)

	if err != nil {
		return models.Engine{}, err
	}
	engine := models.Engine{
		ID:            engineId,
		Displacement:  engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange:      engineReq.CarRange,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
	return engine, nil

}

func (e EngineStore) UpdateEngine(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")
	ctx, span := tracer.Start(ctx, "UpdateEngine-Store")
	defer span.End()
	engineId, err := uuid.Parse(id)
	if err != nil {
		return models.Engine{}, fmt.Errorf("invalid engine id: %v", err)
	}
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("error rolling back transaction: %v", rbErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				fmt.Printf("error committing transaction: %v", commitErr)
			}
		}
	}()
	result, err := tx.ExecContext(ctx,
		"UPDATE engine SET displacement = $2, no_of_cylinders = $3, car_range = $4, updated_at = $5 WHERE id = $1", engineId, engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange, time.Now())
	if err != nil {
		return models.Engine{}, err
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return models.Engine{}, err
	}
	if rowAffected == 0 {
		return models.Engine{}, errors.New("engine not found in database")
	}

	engine := models.Engine{
		ID:            engineId,
		Displacement:  engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange:      engineReq.CarRange,
		UpdatedAt:     time.Now(),
	}
	return engine, nil
}

func (e EngineStore) DeleteEngine(ctx context.Context, id string) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")
	ctx, span := tracer.Start(ctx, "DeleteEngine-Store")
	defer span.End()
	var engine models.Engine
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("error rolling back transaction: %v", rbErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				fmt.Printf("error committing transaction: %v", commitErr)
			}
		}
	}()
	err = tx.QueryRowContext(ctx, "SELECT id, displacement, no_of_cylinders, car_range FROM engine WHERE id = $1", id).Scan(&engine.ID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return engine, errors.New("engine not found in database")
		}
		return engine, err
	}
	result, err := tx.ExecContext(ctx, "DELETE FROM engine WHERE id = $1", id)
	if err != nil {
		return engine, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return engine, err
	}
	if rowsAffected == 0 {
		return engine, errors.New("engine not found in database")
	}
	return engine, nil
}
