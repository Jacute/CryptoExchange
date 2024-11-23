package app_test

import (
	"CryptoExchange/internal/models"
	"CryptoExchange/tests/suite"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLots(t *testing.T) {
	st := suite.New()
	server := httptest.NewServer(st.App.SetupRouter())
	defer server.Close()

	client := &http.Client{}
	lots := GetLots(t, server, client)
	require.Len(t, lots, len(st.Cfg.Lots))
	for i := range lots {
		require.Equal(t, lots[i].Name, st.Cfg.Lots[i])
	}
}

func GetLots(t *testing.T, server *httptest.Server, client *http.Client) []*models.Lot {
	req, err := http.NewRequest(http.MethodGet, server.URL+"/lot", nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var data []*models.Lot
	err = json.Unmarshal(body, &data)
	require.NoError(t, err)

	return data
}
