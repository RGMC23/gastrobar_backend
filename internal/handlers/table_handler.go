package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gastrobar-backend/internal/models"
	"gastrobar-backend/internal/services"

	"github.com/gorilla/mux"
)

type TableHandler struct {
    tableSvc services.TableService
}

func NewTableHandler(tableSvc services.TableService) *TableHandler {
    return &TableHandler{
        tableSvc: tableSvc,
    }
}

func (h *TableHandler) CreateTableHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var table models.Table
        if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        createdTable, err := h.tableSvc.CreateTable(table)
        if err != nil {
            if strings.Contains(err.Error(), "table name cannot be empty") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Table name cannot be empty"})
                return
            }
            if strings.Contains(err.Error(), "table name already exists") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Table name already exists"})
                return
            }
            if strings.Contains(err.Error(), "maximum number of tables reached") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Maximum number of tables reached"})
                return
            }
            log.Printf("Error creating table: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(createdTable)
    }
}

func (h *TableHandler) GetTableHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        tableIDStr := vars["id"]
        tableID, err := strconv.Atoi(tableIDStr)
        if err != nil {
            http.Error(w, "Invalid table ID", http.StatusBadRequest)
            return
        }

        table, err := h.tableSvc.GetTable(tableID)
        if err != nil {
            if strings.Contains(err.Error(), "table not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "Table not found"})
                return
            }
            log.Printf("Error getting table: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(table)
    }
}

func (h *TableHandler) ListTablesHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tables, err := h.tableSvc.ListTables()
        if err != nil {
            log.Printf("Error listing tables: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(tables)
    }
}

func (h *TableHandler) UpdateTableHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        tableIDStr := vars["id"]
        tableID, err := strconv.Atoi(tableIDStr)
        if err != nil {
            http.Error(w, "Invalid table ID", http.StatusBadRequest)
            return
        }

        var table models.Table
        if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        table.ID = tableID

        updatedTable, err := h.tableSvc.UpdateTable(table)
        if err != nil {
            if strings.Contains(err.Error(), "table name cannot be empty") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Table name cannot be empty"})
                return
            }
            if strings.Contains(err.Error(), "table name already exists") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Table name already exists"})
                return
            }
            if strings.Contains(err.Error(), "table not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "Table not found"})
                return
            }
            log.Printf("Error updating table: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(updatedTable)
    }
}