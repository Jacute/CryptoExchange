package jacutesql

import (
	"CryptoExchange/internal/models"
	"CryptoExchange/internal/storage"
	"fmt"
	"strconv"
)

func (s *Storage) GetPairs() ([]*models.Pair, error) {
	const op = "storage.JacuteSQL.GetPairs"

	data, err := s.Query("SELECT pair.pair_pk, pair.first_lot_id, pair.second_lot_id FROM pair")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	pairs := make([]*models.Pair, len(data))
	for i, row := range data {
		id, err := strconv.Atoi(row["pair.pair_pk"])
		if err != nil {
			continue
		}
		buyLotID, err := strconv.Atoi(row["pair.first_lot_id"])
		if err != nil {
			continue
		}
		sellLotID, err := strconv.Atoi(row["pair.second_lot_id"])
		if err != nil {
			continue
		}

		pairs[i] = &models.Pair{
			ID:        id,
			BuyLotID:  buyLotID,
			SellLotID: sellLotID,
		}
	}

	return pairs, nil
}

func (s *Storage) GetPairByID(pairID int) (*models.Pair, error) {
	const op = "storage.JacuteSQL.GetPairs"

	data, err := s.Query("SELECT pair.pair_pk, pair.first_lot_id, pair.second_lot_id FROM pair WHERE pair.pair_pk = '?'", strconv.Itoa(pairID))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrPairNotFound)
	}

	id, err := strconv.Atoi(data[0]["pair.pair_pk"])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	buyLotID, err := strconv.Atoi(data[0]["pair.first_lot_id"])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	sellLotID, err := strconv.Atoi(data[0]["pair.second_lot_id"])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.Pair{
		ID:        id,
		BuyLotID:  buyLotID,
		SellLotID: sellLotID,
	}, nil
}
