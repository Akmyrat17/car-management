package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Car struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Year      string    `json:"year"`
	FuelType  string    `json:"fuel_type"`
	Price     float64   `json:"price"`
	Engine    Engine    `json:"engine"`
	Brand     string    `json:"brand"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CarRequest struct {
	Name     string  `json:"name"`
	Year     string  `json:"year"`
	FuelType string  `json:"fuel_type"`
	Brand    string  `json:"brand"`
	Price    float64 `json:"price"`
	Engine   Engine  `json:"engine"`
}

func CarValidateRequest(carReq CarRequest) error {
	if err := validateName(carReq.Name); err != nil {
		return err
	}

	if err := validateYear(carReq.Year); err != nil {
		return err
	}

	if err := validateFuelType(carReq.FuelType); err != nil {
		return err
	}

	if err := validateBrand(carReq.Brand); err != nil {
		return err
	}
	if err := validateEngine(carReq.Engine); err != nil {
		return err
	}
	if err := validatePrice(carReq.Price); err != nil {
		return err
	}

	return nil
}

func validateName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

func validateYear(year string) error {
	if year == "" {
		return errors.New("year cannot be empty")
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return errors.New("year must be a number")
	}
	currentYear := time.Now().Year()
	if yearInt < 1900 || yearInt > currentYear {
		return errors.New("year must be between 1900 and current year")
	}
	return nil
}

func validateBrand(brand string) error {
	if brand == "" {
		return errors.New("brand cannot be empty")
	}
	return nil
}

func validateFuelType(fuelType string) error {
	validateFuelTypes := []string{"Gasoline", "Diesel", "Electric", "Hybrid"}
	for _, validType := range validateFuelTypes {
		if validType == fuelType {
			return nil
		}
	}
	return errors.New("invalid fuel type selected, please select one of the following: Gasoline, Diesel, Electric, Hybrid")
}

func validateEngine(engine Engine) error {
	if engine.ID == uuid.Nil {
		return errors.New("engine id cannot be empty")
	}
	if engine.Displacement <= 0 {
		return errors.New("engine displacement cannot be empty")
	}
	if engine.NoOfCylinders <= 0 {
		return errors.New("engine number of cylinders cannot be empty")
	}
	if engine.CarRange <= 0 {
		return errors.New("engine car range cannot be empty")
	}
	return nil
}

func validatePrice(price float64) error {
	if price <= 0 {
		return errors.New("price must be greater than 0")
	}
	return nil
}
