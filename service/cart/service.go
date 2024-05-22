package cart

import (
	"fmt"
	"github.com/utsavll0/ecom/types"
)

func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	productIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %v", item.ProductId)
		}
		productIds[i] = item.ProductId
	}
	return productIds, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userId int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}

	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}

	totalPrice := calculateTotalPrice(items, productMap)

	for _, item := range items {
		product := productMap[item.ProductId]
		product.Quantity = product.Quantity - item.Quantity

		err := h.productStore.UpdateProduct(product)
		if err != nil {
			return 0, 0, err
		}
	}

	orderId, err := h.store.CreateOrder(types.Order{
		UserId:  userId,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address",
	})

	if err != nil {
		return 0, 0, err
	}

	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderId:   orderId,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			Price:     float64(productMap[item.ProductId].Price),
		})
	}

	return orderId, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []types.CartItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cartItems is empty")
	}
	for _, cartItem := range cartItems {
		product, ok := products[cartItem.ProductId]
		if !ok {
			return fmt.Errorf("product %d is not available in the store", cartItem.ProductId)
		}

		if product.Quantity <= cartItem.Quantity {
			return fmt.Errorf("product %v is not avaialble in the quantity requested", product.Name)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []types.CartItem, products map[int]types.Product) float64 {
	var total float64
	for _, cartItem := range cartItems {
		product := products[cartItem.ProductId]
		total += float64(product.Price) * float64(cartItem.Quantity)
	}
	return total
}
