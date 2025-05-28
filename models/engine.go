package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Engine struct {
	ID            uuid.UUID `json:"id"`
	Displacement  int64     `json:"displacement"`
	NoOfCylinders int64     `json:"no_of_cylinders"`
	CarRange      int64     `json:"car_range"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type EngineRequest struct {
	Displacement  int64 `json:"displacement"`
	NoOfCylinders int64 `json:"no_of_cylinders"`
	CarRange      int64 `json:"car_range"`
}

func ValidateEngineRequest(engineReq EngineRequest) error {
	if err := validateDisplacement(engineReq.Displacement); err != nil {
		return err
	}
	if err := validateNoOfCylinders(engineReq.NoOfCylinders); err != nil {
		return err
	}
	if err := validateCarRange(engineReq.CarRange); err != nil {
		return err
	}
	return nil
}
func validateDisplacement(displacement int64) error {
	if displacement <= 0 {
		return errors.New("engine displacement cannot be empty, 0 or less")
	}
	return nil
}

func validateNoOfCylinders(noOfCylinders int64) error {
	if noOfCylinders <= 0 {
		return errors.New("engine number of cylinders cannot be empty, 0 or less")
	}
	return nil
}

func validateCarRange(carRange int64) error {
	if carRange <= 0 {
		return errors.New("engine car range cannot be empty, 0 or less")
	}
	return nil
}
