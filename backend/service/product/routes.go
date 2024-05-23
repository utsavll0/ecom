package product

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/utsavll0/ecom/types"
	"github.com/utsavll0/ecom/utils"
	"net/http"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(m *mux.Router) {
	m.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
	m.HandleFunc("/products", h.handleCreateProduct).Methods(http.MethodPost)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		_ = utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusOK, ps)

}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		_ = utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		var e validator.ValidationErrors
		errors.As(err, &e)
		_ = utils.WriteError(w, http.StatusBadRequest, e)
		return
	}

	err := h.store.CreateProduct(types.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
	})
	if err != nil {
		_ = utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusCreated, "Product Created")

}
