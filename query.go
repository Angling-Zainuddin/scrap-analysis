package main

import (
	"fmt"

	"gorm.io/gorm"
)

func GetCustomerByPhone(db *gorm.DB, phone int, phonezero int) (Customers, error) {
	var customer Customers

	err := db.Where(fmt.Sprintf(`CAST(phone AS TEXT) LIKE '%%%d%%' or CAST(whatsapp_number AS TEXT) LIKE '%%%d%%' or CAST(phone AS TEXT) LIKE '%%%d%%' or CAST(whatsapp_number AS TEXT) LIKE '%%%d%%'`, phone, phone, phonezero, phonezero)).Find(&customer).Error
	if err != nil {
		return customer, fmt.Errorf("error when get customer: %s \n", err)
	}

	return customer, nil
}
