package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Akmyrat17/carm/driver"
	carHandler "github.com/Akmyrat17/carm/handler/car"
	engineHandler "github.com/Akmyrat17/carm/handler/engine"
	loginHandler "github.com/Akmyrat17/carm/handler/login"
	"github.com/Akmyrat17/carm/middleware"
	carService "github.com/Akmyrat17/carm/service/car"
	engineService "github.com/Akmyrat17/carm/service/engine"
	carStore "github.com/Akmyrat17/carm/store/car"
	engineStore "github.com/Akmyrat17/carm/store/engine"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Open Telemetry tracing
	traceProvider, err := startTracing()
	if err != nil {
		log.Fatal("Error starting tracing: ", err)
	}

	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shut down tracing, %v", err)
		}
	}()
	otel.SetTracerProvider(traceProvider)

	driver.InitDB()
	defer driver.CloseDB()

	db := driver.GetDB()
	carStore := carStore.New(db)
	carService := carService.NewCarService(carStore)
	carHandler := carHandler.NewCarHandler(carService)

	engineStore := engineStore.New(db)
	engineService := engineService.NewEngineService(engineStore)
	engineHandler := engineHandler.NewEngineHandler(engineService)

	router := mux.NewRouter()
	router.Use(otelmux.Middleware("carm"))
	router.Use(middleware.MetricMiddleware)
	schemaFile := "store/schema.sql"
	if err := executeSchemaFile(db, schemaFile); err != nil {
		log.Fatal("Error executing schema file: ", err)
	}

	router.HandleFunc("/login", loginHandler.LoginHandler).Methods("POST")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/cars/{id}", carHandler.GetCarByID).Methods("GET")
	protected.HandleFunc("/cars", carHandler.CreateCar).Methods("POST")
	protected.HandleFunc("/cars", carHandler.GetCarByBrand).Methods("GET")
	protected.HandleFunc("/cars/{id}", carHandler.UpdateCar).Methods("PUT")
	protected.HandleFunc("/cars/{id}", carHandler.DeleteCar).Methods("DELETE")

	protected.HandleFunc("/engines/{id}", engineHandler.GetEngineById).Methods("GET")
	protected.HandleFunc("/engines", engineHandler.CreateEngine).Methods("POST")
	protected.HandleFunc("/engines/{id}", engineHandler.UpdateEngine).Methods("PUT")
	protected.HandleFunc("/engines/{id}", engineHandler.DeleteEngine).Methods("DELETE")

	router.Handle("/metrics", promhttp.Handler())
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	add := fmt.Sprintf(":%s", port)

	fmt.Println("Starting server on port " + port)
	log.Fatal(http.ListenAndServe(add, router))
}

func executeSchemaFile(db *sql.DB, schemaFile string) error {
	schema, err := os.ReadFile(schemaFile)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		return err
	}
	return nil
}

func startTracing() (*sdktrace.TracerProvider, error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("jaeger:4318"),
			otlptracehttp.WithHeaders(header),
			otlptracehttp.WithInsecure(),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("Error creating new Exporter: %w", err)
	}
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("carm"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay),
		),
		sdktrace.WithResource(res),
	)

	return traceProvider, nil
}
