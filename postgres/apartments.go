package postgres

//import (
//	"fmt"
//	"github.com/jinzhu/gorm"
//	"net/url"
//	"reflect"
//	"savingDeposits"
//	"strconv"
//)
//
//var JsonTagsToFilter = map[string]string{
//	"floor_area_meters":   getJsonTag(savingDeposits.Apartment{}, "FloorAreaMeters"),
//	"price_per_month_usd": getJsonTag(savingDeposits.Apartment{}, "PricePerMonthUsd"),
//	"room_count":          getJsonTag(savingDeposits.Apartment{}, "RoomCount"),
//}
//
//type dbApartmentService struct {
//	Db *gorm.DB
//}
//
//func (ar *dbApartmentService) Create(in savingDeposits.ApartmentCreateInput) (*savingDeposits.ApartmentCreateOutput, error) {
//	ar.Db.Create(&(in.Apartment))
//
//	return &savingDeposits.ApartmentCreateOutput{Apartment: in.Apartment}, nil
//}
//
//func (ar *dbApartmentService) Read(in savingDeposits.ApartmentReadInput) (*savingDeposits.ApartmentReadOutput, error) {
//	apartment, err := getApartment(in.Id, ar.Db)
//	if err != nil {
//		return nil, err
//	}
//
//	return &savingDeposits.ApartmentReadOutput{Apartment: *apartment}, nil
//}
//
//func (ar *dbApartmentService) Find(input savingDeposits.ApartmentFindInput) (*savingDeposits.ApartmentFindOutput, error) {
//	values, err := url.ParseQuery(input.Query)
//	if err != nil {
//		return nil, err
//	}
//
//	tx := ar.Db.New()
//	for dbField, jsonTag := range JsonTagsToFilter {
//		if v, ok := values[jsonTag]; ok {
//			if !ok || len(v) == 0 {
//				continue
//			}
//
//			// TODO potential for injection here
//			tx = tx.Where(fmt.Sprintf("%s = ?", dbField), v[0])
//		}
//	}
//
//	var apartments []savingDeposits.Apartment
//	tx.Find(&apartments)
//	return &savingDeposits.ApartmentFindOutput{Apartments: apartments}, nil
//}
//
//func (ar *dbApartmentService) Update(input savingDeposits.ApartmentUpdateInput) (*savingDeposits.ApartmentUpdateOutput, error) {
//	apartment, err := getApartment(input.Id, ar.Db)
//	if err != nil {
//		return nil, err
//	}
//
//	if err := updateFields(apartment, input.Data); err != nil {
//		return nil, err
//	}
//
//	// Save to DB
//	if err = ar.Db.Save(&apartment).Error; err != nil {
//		return nil, err
//	}
//	return &savingDeposits.ApartmentUpdateOutput{Apartment: *apartment}, nil
//}
//
//func (ar *dbApartmentService) Delete(input savingDeposits.ApartmentDeleteInput) (*savingDeposits.ApartmentDeleteOutput, error) {
//	apartment, err := getApartment(input.Id, ar.Db)
//	if err != nil {
//		return nil, err
//	}
//
//	ar.Db.Delete(&apartment)
//	return &savingDeposits.ApartmentDeleteOutput{Message: "success"}, nil
//}
//
//func getApartment(id string, db *gorm.DB) (*savingDeposits.Apartment, error) {
//	intId, err := strconv.Atoi(id)
//	if err != nil {
//		return nil, err
//	}
//
//	var apartment savingDeposits.Apartment
//	if err = db.First(&apartment, intId).Error; err != nil {
//		if err == gorm.ErrRecordNotFound {
//			return nil, savingDeposits.NotFoundError
//		}
//		return nil, err
//	}
//
//	return &apartment, nil
//}
//
//func updateFields(apartment *savingDeposits.Apartment, data map[string]interface{}) error {
//	if v, ok := data["name"]; ok {
//		apartment.Name = v.(string)
//	}
//
//	if v, ok := data["description"]; ok {
//		apartment.Desc = v.(string)
//	}
//
//	if v, ok := data["floorAreaMeters"]; ok {
//		apartment.FloorAreaMeters = v.(float32)
//	}
//
//	if v, ok := data["pricePerMonthUSD"]; ok {
//		apartment.PricePerMonthUsd = v.(float32)
//	}
//
//	if v, ok := data["roomCount"]; ok {
//		apartment.RoomCount = v.(int)
//	}
//
//	if v, ok := data["latitude"]; ok {
//		apartment.Latitude = v.(float32)
//	}
//
//	if v, ok := data["longitude"]; ok {
//		apartment.Longitude = v.(float32)
//	}
//
//	if v, ok := data["available"]; ok {
//		apartment.Available = v.(bool)
//	}
//
//	return nil
//}
//
//func getJsonTag(v interface{}, fieldName string) string {
//	t := reflect.TypeOf(v)
//	field, ok := t.FieldByName(fieldName)
//	if !ok {
//		return ""
//	}
//
//	return field.Tag.Get("json")
//}
//
//func NewDbApartmentService(db *gorm.DB) *dbApartmentService {
//	return &dbApartmentService{Db: db}
//}
