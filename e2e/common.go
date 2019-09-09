package e2e

import (
	"encoding/json"
	"fmt"
	"log"
	"savingDeposits"
	"savingDeposits/auth"
	"savingDeposits/postgres"
	"savingDeposits/transport"
	"savingDeposits/tst"
	"sync"
	"testing"
)

func getUserToken(t *testing.T, serverUrl, username, pwd string) (string, error) {
	t.Helper()

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, pwd)
	response, err := tst.MakeRequest("POST", serverUrl+"/login", "", []byte(body))
	if err != nil {
		return "", err
	}

	if response.StatusCode >= 399 {
		t.Errorf("Got %d code, expected 2XX", response.StatusCode)
	}

	var token struct {
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&token)

	if err != nil {
		return "", err
	}

	return token.Token, nil
}

func newServer(t *testing.T) (*transport.Server, func()) {
	t.Helper()

	db, err := postgres.ConnectToDB(true)
	tst.Ok(t, err)
	db.AutoMigrate(savingDeposits.DbModels...)

	authN := auth.NewDbAuthnService(db)
	authZ := auth.NewAuthzService()
	depositService := postgres.NewDbSavingDepositService(db)
	usrService := postgres.NewDbUserService(db)

	srv, err := transport.NewServer(db, authN, authZ, depositService,
		usrService)
	tst.Ok(t, err)

	return srv, func() {
		db.DropTableIfExists(savingDeposits.DbModels...)
	}
}

func startServer(wg sync.WaitGroup, addr string, srv *transport.Server) {
	go func() {
		defer wg.Done()
		log.Printf("[ERROR] %s", srv.ServeHTTP(addr))
	}()
}
