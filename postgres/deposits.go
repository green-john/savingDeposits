package postgres

import (
	libErrors "errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"math"
	"net/url"
	"reflect"
	"savingDeposits"
	"strconv"
	"time"
)

var dbNameToJsonField = map[string]string{
	"start_date": getJsonTag(savingDeposits.SavingDeposit{}, "StartDate"),
	"end_date":   getJsonTag(savingDeposits.SavingDeposit{}, "EndDate"),
	"bank_name":  getJsonTag(savingDeposits.SavingDeposit{}, "BankName"),
}

type dbSavingsDepositService struct {
	Db *gorm.DB
}

func turnIntoOneError(errors []error) error {
	finalErr := ""
	for _, err := range errors {
		finalErr += err.Error() + "|"
	}

	return libErrors.New("[dbError]" + finalErr)
}

func (ar *dbSavingsDepositService) Create(input savingDeposits.DepositCreateInput) (*savingDeposits.DepositCreateOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	errors := ar.Db.Create(&(input.SavingDeposit)).GetErrors()
	if len(errors) > 0 {
		return nil, turnIntoOneError(errors)
	}

	user := input.User
	if user.Role != "admin" && uint(user.ID) != input.SavingDeposit.OwnerId {
		return nil, savingDeposits.NotAuthorizedError
	}

	return &savingDeposits.DepositCreateOutput{SavingDeposit: input.SavingDeposit}, nil
}

func (ar *dbSavingsDepositService) Read(input savingDeposits.DepositReadInput) (*savingDeposits.DepositReadOutput, error) {
	deposit, err := getSavingDeposit(input.Id, ar.Db)
	if err != nil {
		return nil, err
	}

	if !canPerformAction(input.User, deposit) {
		return nil, savingDeposits.NotAuthorizedError
	}

	return &savingDeposits.DepositReadOutput{SavingDeposit: *deposit}, nil
}

func (ar *dbSavingsDepositService) Find(input savingDeposits.DepositFindInput) (*savingDeposits.DepositFindOutput, error) {
	values, err := url.ParseQuery(input.Query)
	if err != nil {
		return nil, err
	}

	tx := ar.Db.New()

	for dbField, jsonTag := range dbNameToJsonField {
		if v, ok := values[jsonTag]; ok {
			if !ok || len(v) == 0 {
				continue
			}

			// TODO potential for injection here
			tx = tx.Where(fmt.Sprintf("%s = ?", dbField), v[0])
		}
	}

	if v, ok := values["minAmount"]; ok {
		tx, err = amountGreaterThan(tx, v[0])
		if err != nil {
			return nil, err
		}
	}

	if v, ok := values["maxAmount"]; ok {
		tx, err = amountLessThan(tx, v[0])
		if err != nil {
			return nil, err
		}
	}

	user := input.User
	if user.Role != "admin" {
		tx = tx.Where("owner_id = ?", user.ID)
	}

	var deposits []savingDeposits.SavingDeposit
	tx.Find(&deposits)
	return &savingDeposits.DepositFindOutput{Deposits: deposits}, nil
}

func (ar *dbSavingsDepositService) Update(input savingDeposits.DepositUpdateInput) (*savingDeposits.DepositUpdateOutput, error) {
	deposit, err := getSavingDeposit(input.Id, ar.Db)
	if err != nil {
		return nil, err
	}

	if !canPerformAction(input.User, deposit) {
		return nil, savingDeposits.NotAuthorizedError
	}

	if err := updateFields(deposit, input.Data); err != nil {
		return nil, err
	}

	// Save to DB
	if err = ar.Db.Save(&deposit).Error; err != nil {
		return nil, err
	}
	return &savingDeposits.DepositUpdateOutput{SavingDeposit: *deposit}, nil
}

func (ar *dbSavingsDepositService) Delete(input savingDeposits.DepositDeleteInput) (*savingDeposits.DepositDeleteOutput, error) {
	deposit, err := getSavingDeposit(input.Id, ar.Db)
	if err != nil {
		return nil, err
	}

	if !canPerformAction(input.User, deposit) {
		return nil, savingDeposits.NotAuthorizedError
	}

	ar.Db.Delete(&deposit)
	return &savingDeposits.DepositDeleteOutput{Message: "success"}, nil
}

func (ar *dbSavingsDepositService) GenerateReport(input savingDeposits.GenerateReportInput) (*savingDeposits.GenerateReportOutput, error) {
	values, err := url.ParseQuery(input.Query)
	if err != nil {
		return nil, err
	}

	tx := ar.Db.New()

	if strStartDate, ok := values["startDate"]; ok {
		tx, err = dateGreaterThan(tx, "start_date", strStartDate[0])
		if err != nil {
			return nil, err
		}
	}

	if strEndDate, ok := values["endDate"]; ok {
		tx, err = dateLessThan(tx, "end_date", strEndDate[0])
		if err != nil {
			return nil, err
		}
	}

	user := input.User
	if user.Role != "admin" {
		tx = tx.Where("owner_id = ?", user.ID)
	}

	var deposits []savingDeposits.SavingDeposit
	tx.Find(&deposits)

	response := make([]savingDeposits.ReportEntry, 0, 0)
	for _, d := range deposits {
		revenue := calculateRevenue(d)
		tax := revenue * d.Tax

		if revenue <= 0 {
			tax = 0
		}

		response = append(response, savingDeposits.ReportEntry{
			SavingDeposit: d,
			TotalRevenue:  trimDecimals(revenue),
			TotalTax:      trimDecimals(tax),
			TotalProfit:   trimDecimals(revenue - tax),
		})
	}

	return &savingDeposits.GenerateReportOutput{Deposits: response}, nil
}

// Trims f to 3 decimal places
func trimDecimals(f float64) float64 {
	return math.Round(f*10000) / 10000
}

func calculateRevenue(d savingDeposits.SavingDeposit) float64 {
	startDate := time.Time(d.StartDate)
	endDate := time.Time(d.EndDate)

	totalDays := int(endDate.Sub(startDate).Hours() / 24)
	profitPerDay := d.YearlyInterest / 360

	totalProfit := d.InitialAmount * float64(totalDays) * profitPerDay

	return totalProfit
}

func canPerformAction(user savingDeposits.User, deposit *savingDeposits.SavingDeposit) bool {
	// Admins can do everything
	if user.Role == "admin" {
		return true
	}

	// Otherwise only owners
	return uint(user.ID) == deposit.OwnerId
}

func getSavingDeposit(id string, db *gorm.DB) (*savingDeposits.SavingDeposit, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var SavingDeposit savingDeposits.SavingDeposit
	if err = db.First(&SavingDeposit, intId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, savingDeposits.NotFoundError
		}
		return nil, err
	}

	return &SavingDeposit, nil
}

func updateFields(deposit *savingDeposits.SavingDeposit, data map[string]interface{}) error {
	if v, ok := data["bankName"]; ok {
		deposit.BankName = v.(string)
	}

	if v, ok := data["accountNumber"]; ok {
		deposit.AccountNumber = v.(string)
	}

	if v, ok := data["initialAmount"]; ok {
		deposit.InitialAmount = v.(float64)
	}

	if v, ok := data["yearlyInterest"]; ok {
		deposit.YearlyInterest = v.(float64)
	}

	if v, ok := data["tax"]; ok {
		deposit.Tax = v.(float64)
	}

	if v, ok := data["startDate"]; ok {
		startDate, err := parseDate(v.(string))
		if err != nil {
			return err
		}
		deposit.StartDate = startDate
	}

	if v, ok := data["endDate"]; ok {
		endDate, err := parseDate(v.(string))
		if err != nil {
			return err
		}
		deposit.EndDate = endDate
	}

	if err := deposit.Validate(); err != nil {
		return err
	}

	return nil
}

func parseDate(s string) (savingDeposits.Date, error) {
	t, err := time.Parse(savingDeposits.DateFormat, s)
	if err != nil {
		return savingDeposits.Date{}, err
	}

	return savingDeposits.Date(t), nil
}

func getJsonTag(v interface{}, fieldName string) string {
	t := reflect.TypeOf(v)
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return ""
	}

	return field.Tag.Get("json")
}

func NewDbSavingDepositService(db *gorm.DB) *dbSavingsDepositService {
	return &dbSavingsDepositService{Db: db}
}
