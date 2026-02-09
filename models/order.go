package models

type Order struct {
	Model
	Status string `json:"status" gorm:"type:varchar(20)"`
}
