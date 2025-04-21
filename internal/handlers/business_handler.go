package handlers

import (
    "encoding/json"
    "log"
    "net/http"

    "gastrobar-backend/internal/models"
    "gastrobar-backend/internal/services"

    "github.com/pkg/errors"
)

// BusinessHandler maneja las solicitudes relacionadas con el negocio
type BusinessHandler struct {
    businessSvc services.BusinessService
}

// NewBusinessHandler crea una nueva instancia del manejador de negocio
func NewBusinessHandler(businessSvc services.BusinessService) *BusinessHandler {
    return &BusinessHandler{
        businessSvc: businessSvc,
    }
}

// GetBusinessHandler maneja la solicitud para obtener los datos del negocio
func (h *BusinessHandler) GetBusinessHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        business, err := h.businessSvc.GetBusiness()
        if err != nil {
            if errors.Is(err, errors.Wrap(err, "business not found")) {
                http.Error(w, "Business not found", http.StatusNotFound)
                return
            }
            log.Printf("Error getting business: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(business)
    }
}

// UpdateBusinessHandler maneja la solicitud para actualizar los datos del negocio
func (h *BusinessHandler) UpdateBusinessHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var business models.Business
        if err := json.NewDecoder(r.Body).Decode(&business); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        updatedBusiness, err := h.businessSvc.UpdateBusiness(business)
        if err != nil {
            if errors.Is(err, errors.Wrap(err, "business not found")) {
                http.Error(w, "Business not found", http.StatusNotFound)
                return
            }
            log.Printf("Error updating business: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(updatedBusiness)
    }
}