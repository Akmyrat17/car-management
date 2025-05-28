package car

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Akmyrat17/carm/models"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type CarStore struct {
	db *sql.DB
}

func New(db *sql.DB) *CarStore {
	return &CarStore{db: db}
}

func (c CarStore) GetCarById(ctx context.Context, id string) (models.Car, error) {
	tracer := otel.Tracer("CarStore")
	ctx, span := tracer.Start(ctx, "GetCarById-Store")
	defer span.End()
	var car models.Car

	query := `SELECT c.id,c.name,c.year,c.brand,c.fuel_type,c.price,c.created_at,c.updated_at,e.id,e.displacement,e.no_of_cylinders,e.car_range FROM car c LEFT JOIN engine e ON c.engine_id = e.id WHERE c.id = $1`
	row := c.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType, &car.Price, &car.CreatedAt, &car.UpdatedAt, &car.Engine.ID, &car.Engine.Displacement, &car.Engine.NoOfCylinders, &car.Engine.CarRange)
	if err != nil {
		if err == sql.ErrNoRows {
			return car, nil
		}
		return car, err
	}

	return car, nil
}

func (c CarStore) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	tracer := otel.Tracer("CarStore")
	ctx, span := tracer.Start(ctx, "GetCarByBrand-Store")
	defer span.End()
	var cars []models.Car
	var query string
	if isEngine {
		query = `SELECT c.id,c.name,c.year,c.brand,c.fuel_type,c.price,c.created_at,c.updated_at,e.id,e.displacement,e.no_of_cylinders,e.car_range FROM car c LEFT JOIN engine e ON c.engine_id = e.id WHERE c.brand = $1`
	} else {
		query = `SELECT c.id,c.name,c.year,c.brand,c.fuel_type,c.price,c.created_at,c.updated_at FROM car c WHERE c.brand = $1`
	}

	rows, err := c.db.QueryContext(ctx, query, brand)
	if err != nil {
		return cars, err
	}
	defer rows.Close()
	for rows.Next() {
		var car models.Car
		if isEngine {
			var engine models.Engine
			err := rows.Scan(&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType, &car.Price, &car.CreatedAt, &car.UpdatedAt, &engine.ID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange)
			if err != nil {
				return nil, err
			}
			car.Engine = engine
		} else {
			err := rows.Scan(&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType, &car.Price, &car.CreatedAt, &car.UpdatedAt)
			if err != nil {
				return nil, err
			}
		}
		cars = append(cars, car)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return cars, nil
}

func (c CarStore) CreateCar(ctx context.Context, carReq *models.CarRequest) (models.Car, error) {
	tracer := otel.Tracer("CarStore")
	ctx, span := tracer.Start(ctx, "CreateCar-Store")
	defer span.End()
	var createdCar models.Car
	var engineId uuid.UUID

	err := c.db.QueryRowContext(ctx, "SELECT id FROM engine where id = $1", carReq.Engine.ID).Scan(&engineId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return createdCar, errors.New("engine not found in database")
		}
		return createdCar, err
	}

	carId := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt

	newCar := models.Car{
		ID:        carId,
		Name:      carReq.Name,
		Year:      carReq.Year,
		FuelType:  carReq.FuelType,
		Price:     carReq.Price,
		Engine:    carReq.Engine,
		Brand:     carReq.Brand,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	// Begin the transaction
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return createdCar, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	query := `INSERT INTO car (id, name, year, brand, fuel_type, price, engine_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, name, year, brand, fuel_type, price, created_at, updated_at`
	err = tx.QueryRowContext(ctx, query, newCar.ID, newCar.Name, newCar.Year, newCar.Brand, newCar.FuelType, newCar.Price, newCar.Engine.ID, newCar.CreatedAt, newCar.UpdatedAt).Scan(&createdCar.ID, &createdCar.Name, &createdCar.Year, &createdCar.Brand, &createdCar.FuelType, &createdCar.Price, &createdCar.CreatedAt, &createdCar.UpdatedAt)
	if err != nil {
		return createdCar, err
	}

	return createdCar, nil
}

func (c CarStore) UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (models.Car, error) {
	tracer := otel.Tracer("CarStore")
	ctx, span := tracer.Start(ctx, "UpdateCar-Store")
	defer span.End()
	var updatedCar models.Car

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return updatedCar, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	query :=
		`UPDATE car 
		SET name = $1, year = $2, brand = $3, fuel_type = $4, price = $5, engine_id = $6, updated_at = $7 
			WHERE id = $8 
				RETURNING id, name, year, brand, fuel_type, price, created_at, updated_at`
	err = tx.QueryRowContext(ctx, query, carReq.Name, carReq.Year, carReq.Brand, carReq.FuelType, carReq.Price, carReq.Engine.ID, time.Now(), id).Scan(&updatedCar.ID, &updatedCar.Name, &updatedCar.Year, &updatedCar.Brand, &updatedCar.FuelType, &updatedCar.Price, &updatedCar.CreatedAt, &updatedCar.UpdatedAt)
	if err != nil {
		return updatedCar, err
	}
	return updatedCar, nil
}

func (c CarStore) DeleteCar(ctx context.Context, id string) (models.Car, error) {
	tracer := otel.Tracer("CarStore")
	ctx, span := tracer.Start(ctx, "DeleteCar-Store")
	defer span.End()
	var deletedCar models.Car

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return deletedCar, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.QueryRowContext(ctx, "SELECT id, name, year, brand, fuel_type,engine_id, price, created_at, updated_at FROM car WHERE id = $1", id).Scan(&deletedCar.ID, &deletedCar.Name, &deletedCar.Year, &deletedCar.Brand, &deletedCar.FuelType, &deletedCar.Engine.ID, &deletedCar.Price, &deletedCar.CreatedAt, &deletedCar.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return deletedCar, errors.New("car not found in database")
		}
		return deletedCar, err
	}
	result, err := tx.ExecContext(ctx, "DELETE FROM car WHERE id = $1", id)
	if err != nil {
		return deletedCar, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return deletedCar, err
	}
	if rowsAffected == 0 {
		return deletedCar, errors.New("car not found in database")
	}
	return deletedCar, nil
}
