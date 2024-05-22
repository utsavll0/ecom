package cart

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/utsavll0/ecom/service/auth"
	"github.com/utsavll0/ecom/types"
	"github.com/utsavll0/ecom/utils"
	"net/http"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore) *Handler {
	return &Handler{store: store, productStore: productStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUserIdFromContext(r.Context())
	var cart types.CartCheckoutPayload

	if err := utils.Validate.Struct(cart); err != nil {
		var e validator.ValidationErrors
		errors.As(err, &e)
		_ = utils.WriteJSON(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", e))
		return
	}

	prodcutIds, err := getCartItemsIDs(cart.Items)
	if err != nil {
		_ = utils.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	ps, err := h.productStore.GetProductsById(prodcutIds)
	if err != nil {
		_ = utils.WriteJSON(w, http.StatusInternalServerError, err)
		return
	}

	orderId, totalPrice, err := h.createOrder(ps, cart.Items, userId)

	if err != nil {
		_ = utils.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, map[string]any{
		"orderId":    orderId,
		"totalPrice": totalPrice,
	})
}
