package dao

type User struct {
	Id       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Balance  int    `json:"balance" gorm:"default:1000"`
}
