package jacutesql

import (
	"CryptoExchange/internal/models"
	"CryptoExchange/internal/storage"
	"fmt"
	"strconv"
)

func (s *Storage) SaveOrder(userID, pairID int, quantity, price float64, orderType string) (string, error) {
	const op = "storage.JacuteSQL.SaveOrder"

	userIDStr := strconv.Itoa(userID)
	pairIDStr := strconv.Itoa(pairID)
	quantityStr := strconv.FormatFloat(quantity, 'f', -1, 64)
	priceStr := strconv.FormatFloat(price, 'f', -1, 64)

	id, err := s.Insert(
		"INSERT INTO order VALUES ('?', '?', '?', '?', '?', '?')",
		userIDStr, pairIDStr, quantityStr, priceStr, orderType, "",
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetOrders() ([]*models.Order, error) {
	const op = "storage.JacuteSQL.GetLots"

	data, err := s.Query("SELECT order.order_pk, order.user_id, order.pair_id, order.quantity, order.price, order.type, order.closed FROM order")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	orders := make([]*models.Order, len(data))
	for i, row := range data {
		id, err := strconv.Atoi(row["order.order_pk"])
		if err != nil {
			continue
		}
		userID, err := strconv.Atoi(row["order.user_id"])
		if err != nil {
			continue
		}
		pairID, err := strconv.Atoi(row["order.pair_id"])
		if err != nil {
			continue
		}
		quantity, err := strconv.ParseFloat(row["order.quantity"], 64)
		if err != nil {
			continue
		}
		price, err := strconv.ParseFloat(row["order.price"], 64)
		if err != nil {
			continue
		}

		orders[i] = &models.Order{
			ID:       id,
			UserID:   userID,
			PairID:   pairID,
			Quantity: quantity,
			Price:    price,
			Type:     row["order.type"],
			Closed:   "",
		}
	}

	return orders, nil
}

func (s *Storage) GetOrderByID(orderID string) (*models.Order, error) {
	const op = "storage.JacuteSQL.GetOrderByID"

	id, err := strconv.Atoi(orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid order ID: %w", op, err)
	}

	data, err := s.Query("SELECT order.order_pk, order.user_id, order.pair_id, order.quantity, order.price, order.type, order.closed FROM order WHERE order.order_pk =?", orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(data) == 0 {
		return nil, storage.ErrOrderNotFound
	}

	row := data[0]
	userID, err := strconv.Atoi(row["order.user_id"])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	pairID, err := strconv.Atoi(row["order.pair_id"])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	quantity, err := strconv.ParseFloat(row["order.quantity"], 64)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	price, err := strconv.ParseFloat(row["order.price"], 64)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.Order{
		ID:       id,
		UserID:   userID,
		PairID:   pairID,
		Quantity: quantity,
		Price:    price,
		Type:     row["order.type"],
		Closed:   row["order.closed"],
	}, nil
}

func (s *Storage) DeleteOrder(orderID string) error {
	const op = "storage.JacuteSQL.DeleteOrder"

	err := s.Delete("DELETE FROM order WHERE order.order_pk = ?", orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
