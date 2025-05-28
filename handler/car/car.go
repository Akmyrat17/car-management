package car

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Akmyrat17/carm/models"
	"github.com/Akmyrat17/carm/service"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type CarHandler struct {
	service service.CarServiceInterface
}

func NewCarHandler(service service.CarServiceInterface) *CarHandler {
	return &CarHandler{service: service}
}

func (h *CarHandler) GetCarByID(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "GetCarByd-Handler")
	defer span.End()
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := h.service.GetCarById(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting car by id: ", err)
		return
	}
	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling car: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response: ", err)
		return
	}
}

func (h *CarHandler) GetCarByBrand(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "GetCarByBrand-Handler")
	defer span.End()
	brand := r.URL.Query().Get("brand")

	res, err := h.service.GetCarByBrand(ctx, brand, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting car by brand: ", err)
		return
	}
	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling car: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response: ", err)
		return
	}
}

func (h *CarHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "CreateCar-Handler")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var carReq models.CarRequest
	err = json.Unmarshal(body, &carReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error unmarshalling request body: ", err)
		return
	}

	res, err := h.service.CreateCar(ctx, &carReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating car: ", err)
		return
	}
	body, err = json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling car: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response: ", err)
		return
	}
}

func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "UpdateCar-Handler")
	defer span.End()
	vars := mux.Vars(r)
	id := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var carReq models.CarRequest
	err = json.Unmarshal(body, &carReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error unmarshalling request body: ", err)
		return
	}

	res, err := h.service.UpdateCar(ctx, id, &carReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error updating car: ", err)
		return
	}
	body, err = json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling car: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response: ", err)
		return
	}
}

func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "DeleteCar-Handler")
	defer span.End()
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := h.service.DeleteCar(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error deleting car: ", err)
		return
	}
	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling car: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response: ", err)
		return
	}
}
