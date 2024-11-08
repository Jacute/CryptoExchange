package user_test

import (
	"JacuteCE/tests/suite"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"sync"
	"testing"

	fakeit "github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	cases := []struct {
		Name     string
		Username string
		Status   string
		Err      string
	}{
		{
			Name:     "HappyPath",
			Username: fakeit.Username(),
			Status:   "OK",
			Err:      "",
		},
		{
			Name:     "MaliciousInput",
			Username: ",",
			Status:   "Error",
			Err:      "malicious parameter",
		},
	}

	st := suite.New()
	server := httptest.NewServer(st.App.SetupRouter())
	defer server.Close()

	client := &http.Client{}

	for _, c := range cases {
		t.Run(c.Name, func(tt *testing.T) {
			reqBody := fmt.Sprintf(`{"username": "%s"}`, c.Username)
			req, err := http.NewRequest(http.MethodPost, server.URL+"/user", strings.NewReader(reqBody))
			require.NoError(tt, err)
			req.Header.Set("Content-Type", "application/json")

			res, err := client.Do(req)
			require.NoError(tt, err)
			defer res.Body.Close()
			require.Equal(tt, res.StatusCode, http.StatusOK)

			output, err := io.ReadAll(res.Body)
			require.NoError(tt, err)
			var response map[string]string
			err = json.Unmarshal(output, &response)
			require.NoError(tt, err)

			if c.Err == "" {
				require.Equal(tt, "", response["error"])
				require.Equal(tt, st.Cfg.TokenLen*2, len(response["token"]))
			}

			require.Equal(tt, c.Status, response["status"])
			require.Equal(tt, c.Err, response["error"])
		})
	}
}

func TestUserParallel(t *testing.T) {
	st := suite.New()
	server := httptest.NewServer(st.App.SetupRouter())
	defer server.Close()

	client := &http.Client{}

	gorsCount := runtime.NumCPU()
	runtime.GOMAXPROCS(gorsCount)

	var wg sync.WaitGroup
	wg.Add(gorsCount)
	for i := 0; i < gorsCount; i++ {
		go func() {
			defer wg.Done()
			reqBody := fmt.Sprintf(`{"username": "%s"}`, fakeit.Username())
			req, err := http.NewRequest("POST", server.URL+"/user", strings.NewReader(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			res, err := client.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			output, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			var response map[string]interface{}
			err = json.Unmarshal(output, &response)
			require.NoError(t, err)

			status, ok := response["status"].(string)
			require.True(t, ok)
			if status == "Error" {
				t.Log(response["error"])
			}
			require.Equal(t, "OK", status)
			token, ok := response["token"].(string)
			require.True(t, ok)
			require.Equal(t, st.Cfg.TokenLen*2, len(token))
		}()
	}
	wg.Wait()
}