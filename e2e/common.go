package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"log"
	"savingDeposits"
	"savingDeposits/auth"
	"savingDeposits/postgres"
	"savingDeposits/transport"
	"savingDeposits/tst"
	"sync"
	"testing"
	"time"
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

func newDepositPayload(bankName, accountNumber string, ownerId uint) []byte {
	return []byte(fmt.Sprintf(
		`{
"bankName":"%s",
"accountNumber": "%s",
"initialAmount": 50.0,
"yearlyInterest": 0.5,
"tax": 0.5,
"startDate": "2018-04-20",
"endDate": "2018-04-21",
"ownerId": %d}`, bankName, accountNumber, ownerId))
}

func createNDeposits(t *testing.T, n int, ownerId uint, db *gorm.DB) []uint {
	allIds := make([]uint, 0, 0)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("BNK%d", i)
		endDate := fmt.Sprintf("2019-04-%02d", 10+i)
		id, err := createDeposit(name, float64(i+100), "2019-04-09", endDate, ownerId, db)
		tst.Ok(t, err)

		allIds = append(allIds, id)
	}

	return allIds
}

func parseDate(s string) (savingDeposits.Date, error) {
	t, err := time.Parse(savingDeposits.DateFormat, s)
	if err != nil {
		return savingDeposits.Date{}, err
	}

	return savingDeposits.Date(t), nil

}

func createDeposit(bankName string, initialAmount float64,
	strStartDate string, strEndDate string, ownerId uint, db *gorm.DB) (uint, error) {
	depositsService := postgres.NewDbSavingDepositService(db)

	startDate, err := parseDate(strStartDate)
	if err != nil {
		return 0, err
	}
	endDate, err := parseDate(strEndDate)

	if err != nil {
		return 0, err
	}

	output, err := depositsService.Create(
		savingDeposits.DepositCreateInput{
			SavingDeposit: savingDeposits.SavingDeposit{
				BankName:       bankName,
				AccountNumber:  "no" + bankName,
				InitialAmount:  initialAmount,
				YearlyInterest: 2.25,
				Tax:            .4,
				StartDate:      startDate,
				EndDate:        endDate,
				OwnerId:        ownerId,
			},

			User: savingDeposits.User{Role: "admin"},
		})

	if err != nil {
		return 0, err
	}

	return uint(output.ID), nil
}

func readJson(r io.Reader, v interface{}) error {
	rawContent, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(rawContent, v); err != nil {
		return err
	}
	return nil
}
