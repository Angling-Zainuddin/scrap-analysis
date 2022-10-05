package main

type Customers struct {
	ID       int    `gorm:"column:id"`
	Name     string `gorm:"column:name"`
	Phone    string `gorm:"column:phone"`
	WhatsApp string `gorm:"column:whatsapp_number"`
}
