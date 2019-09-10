package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"savingDeposits"
	"savingDeposits/tst"
	"testing"
	"time"
)

func TestFindDeposits(t *testing.T) {
	// Arrange
	db, err := ConnectToDB(true)
	tst.Ok(t, err)

	db.AutoMigrate(savingDeposits.DbModels...)
	defer db.DropTableIfExists(savingDeposits.DbModels...)

	depositsService := &dbSavingsDepositService{Db: db}

	regularUser, err := addUserToDb("regular", db)
	tst.Ok(t, err)

	err = createDeposits(depositsService, regularUser)
	tst.Ok(t, err)

	allIds := []string{
		"1|b1|2019-04-20|2019-04-23",
		"1|b1|2019-04-23|2019-04-26",
		"1|b1|2019-04-20|2019-04-26",
		"2|b1|2019-04-20|2019-04-23",
		"2|b1|2019-04-23|2019-04-26",
		"2|b1|2019-04-20|2019-04-26",
		"1|b2|2019-04-20|2019-04-23",
		"1|b2|2019-04-23|2019-04-26",
		"1|b2|2019-04-20|2019-04-26",
		"2|b2|2019-04-20|2019-04-23",
		"2|b2|2019-04-23|2019-04-26",
		"2|b2|2019-04-20|2019-04-26",
	}

	start20 := []string{
		"1|b1|2019-04-20|2019-04-23",
		"1|b1|2019-04-20|2019-04-26",
		"2|b1|2019-04-20|2019-04-23",
		"2|b1|2019-04-20|2019-04-26",
		"1|b2|2019-04-20|2019-04-23",
		"1|b2|2019-04-20|2019-04-26",
		"2|b2|2019-04-20|2019-04-23",
		"2|b2|2019-04-20|2019-04-26",
	}

	start20End23 := []string{
		"1|b1|2019-04-20|2019-04-23",
		"2|b1|2019-04-20|2019-04-23",
		"1|b2|2019-04-20|2019-04-23",
		"2|b2|2019-04-20|2019-04-23",
	}

	start20Bank1 := []string{
		"1|b1|2019-04-20|2019-04-23",
		"1|b1|2019-04-20|2019-04-26",
		"2|b1|2019-04-20|2019-04-23",
		"2|b1|2019-04-20|2019-04-26",
	}


	amount2 := []string{
		"2|b1|2019-04-20|2019-04-23",
		"2|b1|2019-04-23|2019-04-26",
		"2|b1|2019-04-20|2019-04-26",
		"2|b2|2019-04-20|2019-04-23",
		"2|b2|2019-04-23|2019-04-26",
		"2|b2|2019-04-20|2019-04-26",
	}

	amount1end26 := []string{
		"1|b1|2019-04-23|2019-04-26",
		"1|b1|2019-04-20|2019-04-26",
		"1|b2|2019-04-23|2019-04-26",
		"1|b2|2019-04-20|2019-04-26",
	}

	for _, elt := range []struct {
		query     string
		resultIds []string
	}{
		{"", allIds},
		{"startDate=2019-04-20", start20},
		{"startDate=2019-04-20&endDate=2019-04-23", start20End23},
		{"startDate=2019-04-20&bankName=b1", start20Bank1},
		{"minAmount=2.0", amount2},
		{"maxAmount=1.0&endDate=2019-04-26", amount1end26},
		{"maxAmount=1.0&minAmount=2.0", []string{}},
		{"minAmount=1.0&maxAmount=2.0", allIds},
	} {
		t.Run(fmt.Sprintf("%s -> %v", elt.query, elt.resultIds), func(t *testing.T) {
			// Act
			res, err := depositsService.Find(savingDeposits.DepositFindInput{
				Query: elt.query,
				User:  regularUser,
			})
			tst.Ok(t, err)

			for idx, deposit := range res.Deposits {
				tst.True(t, deposit.AccountNumber == elt.resultIds[idx],
					fmt.Sprintf("Expected %s, got %s", elt.resultIds[idx], deposit.AccountNumber))
			}
		})
	}
}

// Creates the following apartments:

//  amount |  name  |  startDate  |  endDate  |
//  ===========================================
//     1   |  "b1"  |  2019-4-20  | 2019-4-26
//     1   |  "b1"  |  2019-4-20  | 2019-4-23
//     1   |  "b1"  |  2019-4-23  | 2019-4-26
//     2   |  "b1"  |  2019-4-20  | 2019-4-26
//     2   |  "b1"  |  2019-4-20  | 2019-4-23
//     2   |  "b1"  |  2019-4-23  | 2019-4-26
//     1   |  "b2"  |  2019-4-20  | 2019-4-26
//     1   |  "b2"  |  2019-4-20  | 2019-4-23
//     1   |  "b2"  |  2019-4-23  | 2019-4-26
//     2   |  "b2"  |  2019-4-20  | 2019-4-26
//     2   |  "b2"  |  2019-4-20  | 2019-4-23
//     2   |  "b2"  |  2019-4-23  | 2019-4-26
func createDeposits(s *dbSavingsDepositService, user savingDeposits.User) error {
	allDates := []string{"2019-04-20", "2019-04-23", "2019-04-26"}
	allNames := []string{"b1", "b2"}
	allAmounts := []float64{1, 2}

	for idxN := 0; idxN < len(allNames); idxN++ {
		for idxA := 0; idxA < len(allAmounts); idxA++ {
			payload := newDepositPayload(allAmounts[idxA], allNames[idxN], allDates[0], allDates[1], user)
			_, err := s.Create(payload)
			if err != nil {
				return err
			}
			payload = newDepositPayload(allAmounts[idxA], allNames[idxN], allDates[1], allDates[2], user)
			_, err = s.Create(payload)
			if err != nil {
				return err
			}
			payload = newDepositPayload(allAmounts[idxA], allNames[idxN], allDates[0], allDates[2], user)
			_, err = s.Create(payload)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func newDepositPayload(amount float64, bankName string, rawStartDate, rawEndDate string,
	user savingDeposits.User) savingDeposits.DepositCreateInput {
	accountNumber := fmt.Sprintf("%d|%s|%s|%s",
		int(amount), bankName, rawStartDate, rawEndDate)

	startDate, _ := time.Parse(savingDeposits.DateFormat, rawStartDate)
	endDate, _ := time.Parse(savingDeposits.DateFormat, rawEndDate)

	return savingDeposits.DepositCreateInput{
		SavingDeposit: savingDeposits.SavingDeposit{
			BankName:       bankName,
			AccountNumber:  accountNumber,
			InitialAmount:  amount,
			YearlyTax:      .3,
			YearlyInterest: .3,
			StartDate:      savingDeposits.Date(startDate),
			EndDate:        savingDeposits.Date(endDate),
			OwnerId:        uint(user.ID),
		},
		User: user,
	}
}

func addUserToDb(role string, db *gorm.DB) (savingDeposits.User, error) {
	usrService := NewDbUserService(db)

	res, err := usrService.Create(savingDeposits.UserCreateInput{
		Username: role,
		Password: role,
		Role:     role,
	})

	if err != nil {
		return savingDeposits.User{}, err
	}

	return res.User, nil
}
