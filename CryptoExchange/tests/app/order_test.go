package app_test

import (
	orderdelete "CryptoExchange/internal/http/handlers/order/delete"
	orderpost "CryptoExchange/internal/http/handlers/order/post"
	"CryptoExchange/internal/models"
	"CryptoExchange/tests/suite"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	fakeit "github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestOrderPost_SellAfterBuy(t *testing.T) {
	const orderPairID = 1

	st := suite.New()
	server := httptest.NewServer(st.App.SetupRouter())
	defer server.Close()

	client := &http.Client{}

	res := RegisterUser(t, fakeit.Username(), server, client)
	sellerToken := res["token"]
	require.NotEmpty(t, sellerToken)

	res = RegisterUser(t, fakeit.Username(), server, client)
	buyerToken := res["token"]
	require.NotEmpty(t, buyerToken)

	pairs := Pair(t, client, server)

	var buyLotID, sellLotID int
	for _, pair := range pairs {
		if pair.ID == orderPairID {
			buyLotID = pair.BuyLotID
			sellLotID = pair.SellLotID
			break
		}
	}

	CreateOrder(t, client, server, buyerToken, orderpost.Request{
		PairId:   orderPairID,
		Quantity: 10.0,
		Price:    15.0,
		Type:     "buy",
	})

	balance := Balance(t, client, server, buyerToken)
	for _, lot := range balance {
		if lot.LotID == buyLotID {
			require.Equal(t, 1000.0, lot.Quantity)
		} else if lot.LotID == sellLotID {
			require.Equal(t, 850.0, lot.Quantity)
		}
	}

	CreateOrder(t, client, server, sellerToken, orderpost.Request{
		PairId:   orderPairID,
		Quantity: 5.0,
		Price:    8.0,
		Type:     "sell",
	})

	balance = Balance(t, client, server, buyerToken)
	for _, lot := range balance {
		if lot.LotID == buyLotID {
			require.Equal(t, 1005.0, lot.Quantity)
		} else if lot.LotID == sellLotID {
			require.Equal(t, 885.0, lot.Quantity)
		}
	}

	CreateOrder(t, client, server, sellerToken, orderpost.Request{
		PairId:   orderPairID,
		Quantity: 5.0,
		Price:    8.0,
		Type:     "sell",
	})

	balance = Balance(t, client, server, buyerToken)
	for _, lot := range balance {
		if lot.LotID == buyLotID {
			require.Equal(t, 1010.0, lot.Quantity)
		} else if lot.LotID == sellLotID {
			require.Equal(t, 920.0, lot.Quantity)
		}
	}
}

func TestOrderPost_BuyAfterSell(t *testing.T) {
	const orderPairID = 7

	st := suite.New()
	server := httptest.NewServer(st.App.SetupRouter())
	defer server.Close()

	client := &http.Client{}

	res := RegisterUser(t, fakeit.Username(), server, client)
	sellerToken := res["token"]
	require.NotEmpty(t, sellerToken)

	res = RegisterUser(t, fakeit.Username(), server, client)
	buyerToken := res["token"]
	require.NotEmpty(t, buyerToken)

	pairs := Pair(t, client, server)

	var buyLotID, sellLotID int
	for _, pair := range pairs {
		if pair.ID == orderPairID {
			buyLotID = pair.BuyLotID
			sellLotID = pair.SellLotID
			break
		}
	}

	CreateOrder(t, client, server, sellerToken, orderpost.Request{
		PairId:   orderPairID,
		Quantity: 7.7,
		Price:    22.5,
		Type:     "sell",
	})

	balance := Balance(t, client, server, sellerToken)
	for _, lot := range balance {
		if lot.LotID == buyLotID {
			require.Equal(t, 992.3, lot.Quantity)
		} else if lot.LotID == sellLotID {
			require.Equal(t, 1000.0, lot.Quantity)
		}
	}

	CreateOrder(t, client, server, buyerToken, orderpost.Request{
		PairId:   orderPairID,
		Quantity: 5.0,
		Price:    8.0,
		Type:     "buy",
	})

	balance = Balance(t, client, server, buyerToken)
	for _, lot := range balance {
		if lot.LotID == buyLotID {
			require.Equal(t, 1000.0, lot.Quantity)
		} else if lot.LotID == sellLotID {
			require.Equal(t, 960.0, lot.Quantity)
		}
	}

	CreateOrder(t, client, server, buyerToken, orderpost.Request{
		PairId:   orderPairID,
		Quantity: 5.0,
		Price:    25.0,
		Type:     "buy",
	})

	balance = Balance(t, client, server, buyerToken)
	for _, lot := range balance {
		if lot.LotID == buyLotID {
			require.Equal(t, 1005.0, lot.Quantity)
		} else if lot.LotID == sellLotID {
			require.Equal(t, 847.5, lot.Quantity)
		}
	}

	balance = Balance(t, client, server, sellerToken)
	for _, lot := range balance {
		if lot.LotID == buyLotID {
			require.Equal(t, 992.3, lot.Quantity)
		} else if lot.LotID == sellLotID {
			require.Equal(t, 1112.5, lot.Quantity)
		}
	}
}

func Pair(t *testing.T, client *http.Client, server *httptest.Server) []*models.Pair {
	req, err := http.NewRequest(http.MethodGet, server.URL+"/pair", nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var data []*models.Pair
	err = json.Unmarshal(body, &data)
	require.NoError(t, err)

	require.NotEmpty(t, data)

	return data
}

func TestOrderDelete(t *testing.T) {
	const orderPairID = 7

	st := suite.New()
	server := httptest.NewServer(st.App.SetupRouter())
	defer server.Close()

	client := &http.Client{}

	pairs := Pair(t, client, server)

	var buyLotID, sellLotID int
	for _, pair := range pairs {
		if pair.ID == orderPairID {
			buyLotID = pair.BuyLotID
			sellLotID = pair.SellLotID
			break
		}
	}

	res := RegisterUser(t, fakeit.Username(), server, client)
	sellerToken := res["token"]
	require.NotEmpty(t, sellerToken)

	res = RegisterUser(t, fakeit.Username(), server, client)
	buyerToken := res["token"]
	require.NotEmpty(t, buyerToken)

	CreateOrder(t, client, server, buyerToken, orderpost.Request{
		PairId:   orderPairID,
		Quantity: 5,
		Price:    20,
		Type:     "buy",
	})

	balance := Balance(t, client, server, buyerToken)
	for _, lot := range balance {
		if lot.LotID == buyLotID {
			require.Equal(t, 1000.0, lot.Quantity)
		} else if lot.LotID == sellLotID {
			require.Equal(t, 900.0, lot.Quantity)
		}
	}

	CreateOrder(t, client, server, sellerToken, orderpost.Request{
		PairId:   orderPairID,
		Quantity: 3,
		Price:    18,
		Type:     "sell",
	})

	balance = Balance(t, client, server, buyerToken)
	for _, lot := range balance {
		if lot.LotID == buyLotID {
			require.Equal(t, 1003.0, lot.Quantity)
		} else if lot.LotID == sellLotID {
			require.Equal(t, 906.0, lot.Quantity)
		}
	}

	balance = Balance(t, client, server, sellerToken)
	for _, lot := range balance {
		if lot.LotID == buyLotID {
			require.Equal(t, 997.0, lot.Quantity)
		} else if lot.LotID == sellLotID {
			require.Equal(t, 1054.0, lot.Quantity)
		}
	}

	orders := GetOrder(t, client, server)

	deleteData := DeleteOrder(t, client, server, orders[0].ID, sellerToken)
	require.Equal(t, deleteData["status"], "Error")
	require.Equal(t, deleteData["error"], orderdelete.ErrNotYourOrder.Error)

	deleteData = DeleteOrder(t, client, server, orders[0].ID, buyerToken)
	require.Equal(t, deleteData["status"], "OK")

	newOrders := GetOrder(t, client, server)
	require.Len(t, newOrders, len(orders)-1)
}

func Balance(t *testing.T, client *http.Client, server *httptest.Server, token string) []*models.UserLot {
	req, err := http.NewRequest(http.MethodGet, server.URL+"/balance", nil)
	require.NoError(t, err)
	req.Header.Set("X-USER-TOKEN", token)

	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var data []*models.UserLot
	err = json.Unmarshal(body, &data)
	require.NoError(t, err)

	return data
}

func CreateOrder(t *testing.T, client *http.Client, server *httptest.Server, token string, order orderpost.Request) int {
	reqBody := fmt.Sprintf(`{"pair_id":%d,"quantity":%f,"price":%f,"type":"%s"}`, order.PairId, order.Quantity, order.Price, order.Type)
	req, err := http.NewRequest(http.MethodPost, server.URL+"/order", strings.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-USER-TOKEN", token)

	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	require.NoError(t, err)

	require.Equal(t, data["status"], "OK")
	orderID := int(data["order_id"].(float64))

	return orderID
}

func DeleteOrder(t *testing.T, client *http.Client, server *httptest.Server, id int, token string) map[string]interface{} {
	reqBody := fmt.Sprintf(`{"order_id":%d}`, id)
	req, err := http.NewRequest(http.MethodDelete, server.URL+"/order", strings.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-USER-TOKEN", token)

	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	require.NoError(t, err)

	return data
}

func GetOrder(t *testing.T, client *http.Client, server *httptest.Server) []*models.Order {
	req, err := http.NewRequest(http.MethodGet, server.URL+"/order", nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var data []*models.Order
	err = json.Unmarshal(body, &data)
	require.NoError(t, err)

	return data
}
