package e2e

import (
	"fmt"
	"net/http"
	"savingDeposits/tst"
	"sync"
	"testing"
)

type reportEntry struct {
	ID             uint    `json:"id"`
	BankName       string  `json:"bankName"`
	AccountNumber  string  `json:"accountNumber"`
	InitialAmount  float64 `json:"initialAmount"`
	YearlyInterest float64 `json:"yearlyInterest"`
	YearlyTax      float64 `json:"yearlyTax"`
	StartDate      string  `json:"startDate"`
	EndDate        string  `json:"endDate"`
	TotalRevenue   float64 `json:"totalRevenue"`
	TotalTax       float64 `json:"totalTax"`
	TotalProfit    float64 `json:"totalProfit"`
}

func TestGetReport(t *testing.T) {
	var wg sync.WaitGroup
	const addr = "localhost:8083"
	srv, clean := newServer(t)

	// Delete all the things afterwards
	defer clean()

	serverUrl := fmt.Sprintf("http://%s", addr)
	wg.Add(1)
	startServer(wg, addr, srv)

	regularId, err := createUser("regular", "regular", "regular", srv.Db)
	tst.Ok(t, err)

	_, _ = createDeposit("Bj", 100, "2019-04-20",
		"2019-04-22", regularId, srv.Db)
	_, _ = createDeposit("Bj", 1000, "2019-04-23",
		"2019-04-25", regularId, srv.Db)
	_, _ = createDeposit("Bj", 2000, "2019-04-25",
		"2019-04-27", regularId, srv.Db)

	for _, elt := range []struct {
		query           string
		responseLen     int
		expectedRevenue float64
		expectedTax     float64
		expectedProfit  float64
	}{
		{"startDate=2019-04-20&endDate=2019-04-22", 1, 1.25, 0.5, .75},
		{"startDate=2019-04-20&endDate=2019-04-24", 1, 1.25, 0.5, .75},
	} {
		t.Run(fmt.Sprintf("q:%s", elt.query), func(t *testing.T) {
			token, err := getUserToken(t, serverUrl, "regular", "regular")
			tst.Ok(t, err)
			urlQuery := elt.query
			res, err := tst.MakeRequest("GET", serverUrl+"/report?"+urlQuery,
				token, []byte(""))
			tst.Ok(t, err)

			tst.True(t, res.StatusCode == http.StatusOK,
				fmt.Sprintf("Expected 200, got %d", res.StatusCode))

			var savingDeposits []reportEntry
			err = readJson(res.Body, &savingDeposits)
			tst.True(t, len(savingDeposits) == elt.responseLen, "Expected %d item(s), got %d",
				elt.responseLen, len(savingDeposits))

			tst.True(t, savingDeposits[0].TotalRevenue == elt.expectedRevenue,
				"Expected %f got %f revenue", elt.expectedRevenue, savingDeposits[0].TotalRevenue)

			tst.True(t, savingDeposits[0].TotalTax == elt.expectedTax,
				"Expected %f got %f tax", elt.expectedTax, savingDeposits[0].TotalTax)

			tst.True(t, savingDeposits[0].TotalProfit == elt.expectedProfit,
				"Expected %f got %f profit", elt.expectedProfit, savingDeposits[0].TotalProfit)
		})
	}

}
