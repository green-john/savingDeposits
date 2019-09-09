package e2e

import (
	"encoding/json"
	//"encoding/json"
	"fmt"
	"io/ioutil"

	//"io/ioutil"
	"net/http"
	"savingDeposits/tst"
	"sync"
	"testing"
)

type depositResponse struct {
	ID             uint    `json:"id"`
	BankName       string  `json:"bankName"`
	AccountNumber  string  `json:"accountNumber"`
	InitialAmount  float64 `json:"initialAmount"`
	YearlyInterest float64 `json:"yearlyInterest"`
	YearlyTax      float64 `json:"yearlyTax"`
	StartDate      string  `json:"startDate"`
	EndDate        string  `json:"endDate"`
	OwnerId        uint    `json:"ownerId"`
}

func TestCRUDApartment(t *testing.T) {
	var wg sync.WaitGroup
	const addr = "localhost:8083"
	srv, clean := newServer(t)

	// Delete all the things afterwards
	defer clean()

	serverUrl := fmt.Sprintf("http://%s", addr)
	wg.Add(1)
	startServer(wg, addr, srv)

	_, err := createUser("admin", "admin", "admin", srv.Db)
	tst.Ok(t, err)
	regularId, err := createUser("regular", "regular", "regular", srv.Db)
	tst.Ok(t, err)
	managerId, err := createUser("manager", "manager", "manager", srv.Db)
	tst.Ok(t, err)

	t.Run("CRUD apartment no auth, fail", func(t *testing.T) {
		res, err := tst.MakeRequest("POST", serverUrl+"/deposits", "", []byte(""))
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))

		for _, url := range []string{"/deposits", "/deposits/2"} {
			res, err := tst.MakeRequest("GET", serverUrl+url, "", []byte(""))
			tst.Ok(t, err)

			tst.True(t, res.StatusCode == http.StatusUnauthorized,
				fmt.Sprintf("Expected 401, got %d", res.StatusCode))
		}

		res, err = tst.MakeRequest("PATCH", serverUrl+"/deposits/1", "", []byte(""))
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))

		res, err = tst.MakeRequest("DELETE", serverUrl+"/deposits/1", "", []byte(""))
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))
	})

	t.Run("Create Update Delete apartment with manager, fail", func(t *testing.T) {
		token, err := getUserToken(t, serverUrl, "manager", "manager")
		tst.Ok(t, err)
		newApartmentPayload := newDepositPayload("El Banco", "EB012", managerId)

		res, err := tst.MakeRequest("POST", serverUrl+"/deposits", token, newApartmentPayload)
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusForbidden,
			fmt.Sprintf("Expected 403 got %d", res.StatusCode))

		res, err = tst.MakeRequest("PATCH", serverUrl+"/deposits/1", token, newApartmentPayload)
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusForbidden,
			fmt.Sprintf("Expected 403 got %d", res.StatusCode))

		res, err = tst.MakeRequest("DELETE", serverUrl+"/deposits/1", token, newApartmentPayload)
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusForbidden,
			fmt.Sprintf("Expected 403 got %d", res.StatusCode))
	})

	t.Run("CRUD apartment admin, success", func(t *testing.T) {
		// Get user token
		token, err := getUserToken(t, serverUrl, "admin", "admin")
		tst.Ok(t, err)

		// Create
		payload := newDepositPayload("Bank", "AC01", regularId)
		res, err := tst.MakeRequest("POST", serverUrl+"/deposits", token, payload)
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusCreated, fmt.Sprintf("Expected 201 got %d", res.StatusCode))
		rawContent, err := ioutil.ReadAll(res.Body)
		tst.Ok(t, err)

		var createdDeposit depositResponse
		err = json.Unmarshal(rawContent, &createdDeposit)

		tst.True(t, createdDeposit.ID >= 1, "Expected id greater than 0")
		tst.True(t, createdDeposit.BankName == "Bank", "Got different bank name")
		tst.True(t, createdDeposit.AccountNumber == "AC01", "Got different account number")
		tst.True(t, createdDeposit.StartDate == "2018-04-20", "Got unexpected start date %s",
			createdDeposit.StartDate)
		tst.True(t, createdDeposit.EndDate == "2018-04-21", "Got unexpected end date %s",
			createdDeposit.EndDate)
		tst.True(t, createdDeposit.OwnerId == regularId, "Got different owner id")

		// Read
		depositUrl := fmt.Sprintf("%s/deposits/%d", serverUrl, createdDeposit.ID)
		res, err = tst.MakeRequest("GET", depositUrl, token, []byte(""))
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusOK,
			fmt.Sprintf("Expected 200, got %d", res.StatusCode))

		rawContent, err = ioutil.ReadAll(res.Body)
		tst.Ok(t, err)

		var retrievedDeposit depositResponse
		err = json.Unmarshal(rawContent, &retrievedDeposit)
		tst.Ok(t, err)

		tst.True(t, retrievedDeposit.ID == createdDeposit.ID, fmt.Sprintf("Unexpected id %d",
			retrievedDeposit.ID))

		// Update
		newData := []byte(`{"id": 200, "bankName": "New Bank", "yearlyInterest": 0.3, "endDate":"2019-04-20"}`)
		res, err = tst.MakeRequest("PATCH", depositUrl, token, newData)
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusOK,
			fmt.Sprintf("Expected 200, got %d", res.StatusCode))

		res, err = tst.MakeRequest("GET", depositUrl, token, []byte(""))
		tst.Ok(t, err)
		rawContent, err = ioutil.ReadAll(res.Body)
		tst.Ok(t, err)

		var updatedDeposit depositResponse
		err = json.Unmarshal(rawContent, &updatedDeposit)
		tst.Ok(t, err)
		tst.True(t, updatedDeposit.ID == retrievedDeposit.ID,
			"Expected id to be %d, got %d", updatedDeposit.ID, retrievedDeposit.ID)
		tst.True(t, updatedDeposit.BankName == "New Bank",
			"Expected bank name to be NewBank, got %s", updatedDeposit.BankName)
		tst.True(t, updatedDeposit.YearlyInterest == 0.3,
			"Expected interest to be 0.3, got %f", updatedDeposit.YearlyInterest)
		tst.True(t, updatedDeposit.YearlyTax == createdDeposit.YearlyTax,
			"Expected tax to be %f, got %f", createdDeposit.YearlyTax, updatedDeposit.YearlyTax)
		tst.True(t, updatedDeposit.StartDate == createdDeposit.StartDate,
			"Expected start date to be %s, got %s", createdDeposit.StartDate, updatedDeposit.StartDate)
		tst.True(t, updatedDeposit.EndDate == "2019-04-20",
			"Expected end date to be %s, got %s", createdDeposit.EndDate, updatedDeposit.EndDate)

		// Delete
		res, err = tst.MakeRequest("DELETE", depositUrl, token, []byte(""))
		tst.Ok(t, err)
		tst.True(t, res.StatusCode == http.StatusNoContent,
			fmt.Sprintf("Expected 204, got %d", res.StatusCode))

		res, err = tst.MakeRequest("GET", depositUrl, token, []byte(""))
		tst.Ok(t, err)
		tst.True(t, res.StatusCode == http.StatusNotFound, "Found deleted shite")
	})
}

//func TestReadAllApartmentsAndSearch(t *testing.T) {
//	var wg sync.WaitGroup
//	const addr = "localhost:8083"
//	srv, clean := newServer(t)
//	defer clean()
//
//	serverUrl := fmt.Sprintf("http://%s", addr)
//
//	wg.Add(1)
//	startServer(wg, addr, srv)
//
//	_, err := createUser("admin", "admin", "admin", srv.Db)
//	tst.Ok(t, err)
//	realtorId, err := createUser("realtor", "realtor", "realtor", srv.Db)
//	tst.Ok(t, err)
//	_, err = createUser("client", "client", "client", srv.Db)
//	tst.Ok(t, err)
//
//	create10Apartments(t, realtorId, srv.Db)
//
//	t.Run("Read all apartments client, realtor, admin, success", func(t *testing.T) {
//		for _, user := range []string{"client", "realtor", "admin"} {
//			token, err := getUserToken(t, serverUrl, user, user)
//			tst.Ok(t, err)
//
//			// Act
//			res, err := tst.MakeRequest("GET", serverUrl+"/apartments", token, []byte(""))
//			tst.Ok(t, err)
//
//			// True
//			tst.True(t, res.StatusCode == http.StatusOK,
//				fmt.Sprintf("Expected 200, got %d", res.StatusCode))
//
//			var returnedApartments []apartmentResponse
//			decoder := json.NewDecoder(res.Body)
//			err = decoder.Decode(&returnedApartments)
//			tst.Ok(t, err)
//
//			tst.True(t, len(returnedApartments) == 10,
//				fmt.Sprintf("Expected 10 apartments, got %d", len(returnedApartments)))
//		}
//	})
//
//	t.Run("Search apartments by room count", func(t *testing.T) {
//		token, err := getUserToken(t, serverUrl, "client", "client")
//		tst.Ok(t, err)
//
//		// Act
//		res, err := tst.MakeRequest("GET", serverUrl+"/apartments?roomCount=4", token, []byte(""))
//		tst.Ok(t, err)
//
//		// True
//		tst.True(t, res.StatusCode == http.StatusOK,
//			fmt.Sprintf("Expected 200, got %d", res.StatusCode))
//
//		var returnedApartments []apartmentResponse
//		decoder := json.NewDecoder(res.Body)
//		err = decoder.Decode(&returnedApartments)
//		tst.Ok(t, err)
//
//		tst.True(t, len(returnedApartments) == 5,
//			fmt.Sprintf("Expected 5 apartments, got %d", len(returnedApartments)))
//	})
//}
//
func newDepositPayload(bankName, accountNumber string, ownerId uint) []byte {
	return []byte(fmt.Sprintf(
		`{
"bankName":"%s",
"accountNumber": "%s",
"initialAmount": 50.0,
"yearlyInterest": 0.5,
"yearlyTax": 0.5,
"startDate": "2018-04-20",
"endDate": "2018-04-21",
"ownerId": %d}`, bankName, accountNumber, ownerId))
}

//
//func create10Apartments(t *testing.T, realtorId uint, db *gorm.DB) {
//	for i := 0; i < 5; i++ {
//		name := fmt.Sprintf("apt%d", i)
//		desc := fmt.Sprintf("desc%d", i)
//		_, err := createApartment(name, desc, 2, realtorId, db)
//		tst.Ok(t, err)
//	}
//
//	for i := 0; i < 5; i++ {
//		name := fmt.Sprintf("apt%d", 5+i)
//		desc := fmt.Sprintf("desc%d", 5+i)
//		_, err := createApartment(name, desc, 4, realtorId, db)
//		tst.Ok(t, err)
//	}
//}
//
//func createApartment(name, desc string, roomCount int, realtorId uint, db *gorm.DB) (uint, error) {
//	apartmentResource := postgres.NewDbApartmentService(db)
//
//	output, err := apartmentResource.Create(
//		savingDeposits.ApartmentCreateInput{
//			Apartment: savingDeposits.Apartment{
//				Name:             name,
//				Desc:             desc,
//				RoomCount:        roomCount,
//				PricePerMonthUsd: 500.0,
//				Latitude:         41.761536,
//				Longitude:        12.315237,
//				RealtorId:        realtorId,
//				Available:        true,
//			},
//		})
//	if err != nil {
//		return 0, err
//	}
//
//	return uint(output.ID), nil
//}
