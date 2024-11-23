package jacutesql

import (
	"CryptoExchange/internal/models"
	"CryptoExchange/internal/storage"
	"fmt"
	"strconv"
	"time"
)

func (s *Storage) SaveOrder(order *models.Order) (string, error) {
	const op = "storage.JacuteSQL.SaveOrder"

	userIDStr := strconv.Itoa(order.UserID)
	pairIDStr := strconv.Itoa(order.PairID)
	quantityStr := strconv.FormatFloat(order.Quantity, 'f', -1, 64)
	priceStr := strconv.FormatFloat(order.Price, 'f', -1, 64)

	id, err := s.Insert(
		"INSERT INTO order VALUES ('?', '?', '?', '?', '?', '?')",
		userIDStr, pairIDStr, quantityStr, priceStr, order.Type, order.Closed,
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
			Closed:   row["order.closed"],
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

	data, err := s.Query("SELECT order.order_pk, order.user_id, order.pair_id, order.quantity, order.price, order.type, order.closed FROM order WHERE order.order_pk = ?", orderID)
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

func (s *Storage) GetOrderForOperation(price float64, pairID, orderType string) (*models.Order, error) {
	const op = "storage.JacuteSQL.CheckForBuying"

	data, err := s.Query("SELECT order.order_pk, order.user_id, order.pair_id, order.quantity, order.price, order.type, order.closed FROM order WHERE order.type = '?' AND order.pair_id = '?'", orderType, pairID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, row := range data {
		if val, ok := row["order.closed"]; ok && val != "0" {
			continue
		}

		id, err := strconv.Atoi(row["order.order_pk"])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

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
		orderPrice, err := strconv.ParseFloat(row["order.price"], 64)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		order := &models.Order{
			ID:       id,
			UserID:   userID,
			PairID:   pairID,
			Quantity: quantity,
			Price:    orderPrice,
			Type:     row["order.type"],
			Closed:   row["order.closed"],
		}
		fmt.Println(order, price)
		if orderType == "buy" {
			if order.Price >= price {
				return order, nil
			}
		} else {
			if order.Price <= price {
				return order, nil
			}
		}
	}
	return nil, storage.ErrOrderNotFound
}

func (s *Storage) Buy(buyerOrder *models.Order, sellerOrder *models.Order) error {
	const op = "storage.JacuteSQL.Buy"

	pair, err := s.GetPairByID(buyerOrder.PairID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	curTime := strconv.Itoa(int(time.Now().Unix()))
	buyerOrder.Closed = curTime
	sellerOrder.Closed = curTime

	err = s.Delete("DELETE FROM order WHERE order.order_pk = '?' OR order.order_pk = '?'", strconv.Itoa(buyerOrder.ID), strconv.Itoa(sellerOrder.ID))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	buyerID := strconv.Itoa(buyerOrder.UserID)
	buyLotID := strconv.Itoa(pair.BuyLotID)
	sellerID := strconv.Itoa(sellerOrder.UserID)
	sellLotID := strconv.Itoa(pair.SellLotID)

	fmt.Println(1, buyerOrder.Quantity, sellerOrder.Quantity)
	if buyerOrder.Quantity > sellerOrder.Quantity {
		// deposit lot to the buyer account
		_, err = s.AddMoney(buyerID, buyLotID, sellerOrder.Quantity)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if sellerOrder.Price < buyerOrder.Price {
			// deposit remains to the buyer account
			remains := (buyerOrder.Price - sellerOrder.Price) * sellerOrder.Quantity
			_, err = s.AddMoney(buyerID, sellLotID, remains)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
		s.AddMoney(sellerID, sellLotID, sellerOrder.Quantity*sellerOrder.Price)
		s.Pay(sellerID, buyLotID, sellerOrder.Quantity)

		diff := buyerOrder.Quantity - sellerOrder.Quantity

		_, err := s.SaveOrder(&models.Order{
			UserID:   buyerOrder.UserID,
			PairID:   buyerOrder.PairID,
			Quantity: diff,
			Price:    buyerOrder.Price,
			Type:     "buy",
			Closed:   "0",
		})
		if err != nil {
			return fmt.Errorf("%s: can't create new order for buyer: %w", op, err)
		}
		// close seller order
		_, err = s.SaveOrder(sellerOrder)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else if sellerOrder.Quantity >= buyerOrder.Quantity {
		// deposit lot to the buyer account
		_, err = s.AddMoney(buyerID, buyLotID, buyerOrder.Quantity)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if sellerOrder.Price < buyerOrder.Price {
			// deposit remains to the buyer account
			remains := (buyerOrder.Price - sellerOrder.Price) * buyerOrder.Quantity
			_, err = s.AddMoney(buyerID, sellLotID, remains)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
		s.AddMoney(sellerID, sellLotID, buyerOrder.Quantity*sellerOrder.Price)
		s.Pay(sellerID, buyLotID, buyerOrder.Quantity)

		if sellerOrder.Quantity > buyerOrder.Quantity {
			diff := sellerOrder.Quantity - buyerOrder.Quantity

			_, err := s.SaveOrder(&models.Order{
				UserID:   sellerOrder.UserID,
				PairID:   sellerOrder.PairID,
				Quantity: diff,
				Price:    sellerOrder.Price,
				Type:     "sell",
				Closed:   "0",
			})
			if err != nil {
				return fmt.Errorf("%s: can't create new order for seller: %w", op, err)
			}

			// close buyer order
			_, err = s.SaveOrder(buyerOrder)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		} else {
			// close orders
			_, err = s.SaveOrder(sellerOrder)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			_, err = s.SaveOrder(buyerOrder)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	}
	return nil
}
