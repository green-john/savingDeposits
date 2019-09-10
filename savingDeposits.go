package savingDeposits

import (
	"database/sql/driver"
	"errors"
	"math"
	"strings"
	"time"
)

type Date time.Time

const DateFormat = "2006-01-02"

type SavingDeposit struct {
	// Primary key
	ID uid `gorm:"primary_key" json:"id"`

	// Bank name
	BankName string `json:"bankName"`

	// Account number
	AccountNumber string `json:"accountNumber"`

	// Initial Amount
	InitialAmount float64 `json:"initialAmount"`

	// Yearly interest
	YearlyInterest float64 `json:"yearlyInterest"`

	// Yearly taxes
	YearlyTax float64 `json:"yearlyTax"`

	// Date added
	StartDate Date `gorm:"type:date" json:"startDate"`

	// Date added
	EndDate Date `gorm:"type:date" json:"endDate"`

	// Realtor associated with this SavingDeposit
	Owner   User `json:"-" gorm:"foreignkey:OwnerId"`
	OwnerId uint `json:"ownerId"`
}

func (uid) UnmarshalJSON([]byte) error {
	return nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	println(s)

	t, err := time.Parse(DateFormat, s)
	if err != nil {
		return err
	}

	*d = Date(t)
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return []byte("\"" + t.Format(DateFormat) + "\""), nil
}

func (d Date) Value() (driver.Value, error) {
	return time.Time(d), nil
}

func (d *Date) Scan(value interface{}) error {
	*d = Date(value.(time.Time))
	return nil
}

// Validates data for a new deposits.
func (s *SavingDeposit) Validate() error {
	allErrors := ""

	if strings.Trim(s.BankName, " ") == "" {
		allErrors += "Bank name can't be empty\n"
	}

	if strings.Trim(s.AccountNumber, " ") == "" {
		allErrors += "Account number can't be empty\n"
	}

	if s.InitialAmount <= 0 {
		allErrors += "Initial amount must be greater than 0\n"
	}

	if !isFraction(s.YearlyTax) {
		allErrors += "Yearly tax should be positive and between [0.0, 1.0]"
	}

	if !isFraction(math.Abs(s.YearlyInterest)) {
		allErrors += "Yearly interest should be between [-1.0, 1.0]"
	}

	startDate := time.Time(s.StartDate)
	endDate := time.Time(s.EndDate)

	if !startDate.Before(endDate) {
		allErrors += "Start date must be before end date\n"
	}

	// TODO validate owner is not zero

	if allErrors == "" {
		return nil
	}

	return errors.New("\n" + allErrors)
}

func isFraction(n float64) bool {
	return n >= 0 && n <= 1.0
}

type DespositCreateInput struct {
	SavingDeposit

	// Auth user
	User User
}

type DepositCreateOutput struct {
	SavingDeposit
}

type DepositReadInput struct {
	// ID to lookup the SavingDeposit
	Id string

	// Auth user
	User User
}

type DepositReadOutput struct {
	SavingDeposit
}

type DepositFindInput struct {
	Query string

	// Auth user
	User User
}

type DepositFindOutput struct {
	Deposits []SavingDeposit
}

func (o *DepositFindOutput) Public() interface{} {
	return o.Deposits
}

type DepositUpdateInput struct {
	Id   string
	Data map[string]interface{}

	// Auth user
	User User
}

type DepositUpdateOutput struct {
	SavingDeposit
}

type DepositDeleteInput struct {
	Id string

	// Auth user
	User User
}

type DepositDeleteOutput struct {
	Message string
}

type DepositsService interface {
	Create(DespositCreateInput) (*DepositCreateOutput, error)
	Read(DepositReadInput) (*DepositReadOutput, error)
	Find(DepositFindInput) (*DepositFindOutput, error)
	Update(DepositUpdateInput) (*DepositUpdateOutput, error)
	Delete(DepositDeleteInput) (*DepositDeleteOutput, error)
}
