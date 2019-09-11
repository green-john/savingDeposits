package postgres

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

func amountGreaterThan(tx *gorm.DB, strAmount string) (*gorm.DB, error) {
	amount, err := strconv.ParseFloat(strAmount, 64)
	if err != nil {
		return nil, err
	}

	return tx.Where("initial_amount >= ?", amount), nil
}

func amountLessThan(tx *gorm.DB, strAmount string) (*gorm.DB, error) {
	amount, err := strconv.ParseFloat(strAmount, 64)
	if err != nil {
		return nil, err
	}

	return tx.Where("initial_amount <= ?", amount), nil
}

func dateGreaterThan(tx *gorm.DB, dateField, strDate string) (*gorm.DB, error) {
	date, err := parseDate(strDate)
	if err != nil {
		return nil, err
	}

	query := dateField + " >= ?"
	return tx.Where(query, date), nil
}

func dateLessThan(tx *gorm.DB, dateField, strDate string) (*gorm.DB, error) {
	date, err := parseDate(strDate)
	if err != nil {
		return nil, err
	}

	query := dateField + " <= ?"
	return tx.Where(query, date), nil
}
