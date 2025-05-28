package engine

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

type EngineHandler struct {
	engineService service.EngineServiceInterface
}

func NewEngineHandler(engineService service.EngineServiceInterface) *EngineHandler {
	return &EngineHandler{engineService: engineService}
}

func (e EngineHandler) GetEngineById(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(r.Context(), "GetEngineById-Handler")
	defer span.End()
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := e.engineService.GetEngineById(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting engine by id: ", err)
		return
	}
	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling engine: ", err)
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

func (e EngineHandler) CreateEngine(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(r.Context(), "CreaTEeNGINE-Handler")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var engineReq models.EngineRequest
	err = json.Unmarshal(body, &engineReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error unmarshalling request body: ", err)
		return
	}

	res, err := e.engineService.CreateEngine(ctx, &engineReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating engine: ", err)
		return
	}
	body, err = json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling engine: ", err)
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

func (e EngineHandler) UpdateEngine(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(r.Context(), "UpdateEngine-Handler")
	defer span.End()
	vars := mux.Vars(r)
	id := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var engineReq models.EngineRequest
	err = json.Unmarshal(body, &engineReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error unmarshalling request body: ", err)
		return
	}

	res, err := e.engineService.UpdateEngine(ctx, id, &engineReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error updating engine: ", err)
		return
	}
	body, err = json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling engine: ", err)
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

func (e EngineHandler) DeleteEngine(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(r.Context(), "DeleteEngine-Handler")
	defer span.End()
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := e.engineService.DeleteEngine(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error deleting engine: ", err)
		return
	}
	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling engine: ", err)
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
