package handler

import (
	"app/internal"
	"app/platform/web/request"
	"app/platform/web/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// NewHandlerWarehouse creates a new handler for Warehouse.
func NewHandlerWarehouse(rp internal.RepositoryWarehouse) (h *HandlerWarehouse) {
	h = &HandlerWarehouse{
		rp: rp,
	}
	return
}

// HandlerWarehouse is a handler for Warehouse.
type HandlerWarehouse struct {
	// rp is the repository for Warehouse.
	rp internal.RepositoryWarehouse
}

// ProductJSON is a Warehouse in JSON format.
type WarehouseJSONResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
	Capacity  int    `json:"capacity"`
}

type WarehouseJSONRequest struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
	Capacity  int    `json:"capacity"`
}

// GetById gets a string by id.
func (h *HandlerWarehouse) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// - path parameter: id
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.JSON(w, http.StatusBadRequest, "invalid id")
			return
		}

		// process
		// - find Warehouse by id
		p, err := h.rp.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrRepositoryWarehouseNotFound):
				response.JSON(w, http.StatusNotFound, "Warehouse not found")
			default:
				response.JSON(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// - serialize Warehouse to JSON
		data := WarehouseJSONResponse{
			Id:        p.Id,
			Name:      p.Name,
			Address:   p.Address,
			Telephone: p.Telephone,
			Capacity:  p.Capacity,
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

// Create creates a Warehouse.
func (h *HandlerWarehouse) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// - body
		var body WarehouseJSONRequest
		err := request.JSON(r, &body)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, "invalid body")
			return
		}

		// process
		// - save Warehouse

		warehouse := internal.Warehouse{
			WarehouseAttributes: internal.WarehouseAttributes{
				Name:      body.Name,
				Address:   body.Address,
				Telephone: body.Telephone,
				Capacity:  body.Capacity,
			},
		}

		err = h.rp.Save(&warehouse)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// response
		// - serialize product to JSON
		data := WarehouseJSONResponse{
			Id:        warehouse.Id,
			Name:      warehouse.Name,
			Address:   warehouse.Address,
			Telephone: warehouse.Telephone,
			Capacity:  warehouse.Capacity,
		}
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

// ReportProducts gets the number of products by warehouse.
func (h *HandlerWarehouse) ReportProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// - query parameter: id
		idParam := r.URL.Query().Get("id")

		var id int
		if idParam != "" {
			var err error
			id, err = strconv.Atoi(idParam)
			if err != nil {
				response.JSON(w, http.StatusBadRequest, "invalid id")
				return
			}
		}

		// process
		// - find Warehouse by id
		warehouses, err := h.rp.ReportProducts(id)

		if err != nil {
			switch {
			default:
				response.JSON(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    warehouses,
		})
	}
}

func (h *HandlerWarehouse) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warehouses, err := h.rp.GetAll()

		if err != nil {
			response.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

		response.JSON(w, http.StatusOK, warehouses)
	}
}
