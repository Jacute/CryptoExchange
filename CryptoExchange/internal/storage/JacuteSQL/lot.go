package jacutesql

import (
	"CryptoExchange/internal/models"
	"fmt"
	"strconv"
)

func (s *Storage) GetLots() ([]*models.Lot, error) {
	const op = "storage.JacuteSQL.GetLots"

	data, err := s.Query("SELECT lot.lot_pk, lot.name FROM lot")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	lots := make([]*models.Lot, len(data))
	for i, row := range data {
		id, err := strconv.Atoi(row["lot.lot_pk"])
		if err != nil {
			continue
		}

		lots[i] = &models.Lot{
			ID:   id,
			Name: row["lot.name"],
		}
	}

	return lots, nil
}

func (s *Storage) AddLots(userID string, quantity string) error {
	const op = "storage.JacuteSQL.SaveLot"

	data, err := s.Query("SELECT lot.lot_pk FROM lot")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, row := range data {
		_, err := s.Insert("INSERT INTO user_lot VALUES ('?', '?', '?')", userID, row["lot.lot_pk"], quantity)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
