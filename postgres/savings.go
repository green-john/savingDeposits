package postgres

import (
	libErrors "errors"
	"github.com/jinzhu/gorm"
	"reflect"
	"savingDeposits"
	"strconv"
	"time"
)

//var JsonTagsToFilter = map[string]string{
//	"floor_area_meters":   getJsonTag(savingDeposits.SavingDeposit{}, "FloorAreaMeters"),
//	"price_per_month_usd": getJsonTag(savingDeposits.SavingDeposit{}, "PricePerMonthUsd"),
//	"room_count":          getJsonTag(savingDeposits.SavingDeposit{}, "RoomCount"),
//}

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

func (ar *dbSavingsDepositService) Create(in savingDeposits.DespositCreateInput) (*savingDeposits.DepositCreateOutput, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}
	errors := ar.Db.Create(&(in.SavingDeposit)).GetErrors()
	if len(errors) > 0 {
		return nil, turnIntoOneError(errors)
	}

	return &savingDeposits.DepositCreateOutput{SavingDeposit: in.SavingDeposit}, nil
}

func (ar *dbSavingsDepositService) Read(in savingDeposits.DepositReadInput) (*savingDeposits.DepositReadOutput, error) {
	SavingDeposit, err := getSavingDeposit(in.Id, ar.Db)
	if err != nil {
		return nil, err
	}

	return &savingDeposits.DepositReadOutput{SavingDeposit: *SavingDeposit}, nil
}

func (ar *dbSavingsDepositService) Find(input savingDeposits.DepositFindInput) (*savingDeposits.DepositFindOutput, error) {
	//values, err := url.ParseQuery(input.Query)
	//if err != nil {
	//	return nil, err
	//}

	//tx := ar.Db.New()
	//for dbField, jsonTag := range JsonTagsToFilter {
	//	if v, ok := values[jsonTag]; ok {
	//		if !ok || len(v) == 0 {
	//			continue
	//		}
	//
	//		// TODO potential for injection here
	//		tx = tx.Where(fmt.Sprintf("%s = ?", dbField), v[0])
	//	}
	//}

	//var SavingDeposits []savingDeposits.SavingDeposit
	//tx.Find(&SavingDeposits)
	return &savingDeposits.DepositFindOutput{}, nil
}

func (ar *dbSavingsDepositService) Update(input savingDeposits.DepositUpdateInput) (*savingDeposits.DepositUpdateOutput, error) {
	SavingDeposit, err := getSavingDeposit(input.Id, ar.Db)
	if err != nil {
		return nil, err
	}

	if err := updateFields(SavingDeposit, input.Data); err != nil {
		return nil, err
	}

	// Save to DB
	if err = ar.Db.Save(&SavingDeposit).Error; err != nil {
		return nil, err
	}
	return &savingDeposits.DepositUpdateOutput{SavingDeposit: *SavingDeposit}, nil
}

func (ar *dbSavingsDepositService) Delete(input savingDeposits.DepositDeleteInput) (*savingDeposits.DepositDeleteOutput, error) {
	SavingDeposit, err := getSavingDeposit(input.Id, ar.Db)
	if err != nil {
		return nil, err
	}

	ar.Db.Delete(&SavingDeposit)
	return &savingDeposits.DepositDeleteOutput{Message: "success"}, nil
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

	if v, ok := data["yearlyTax"]; ok {
		deposit.YearlyTax = v.(float64)
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
