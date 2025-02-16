package dao

type Transaction struct {
	Id         uint `json:"id" gorm:"primaryKey"`
	FromUserId uint `json:"fromUserId"`
	ToUserId   uint `json:"toUserId"`
	Amount     int  `json:"amount"`
	FromUser   User `gorm:"foreignKey:FromUserId;references:Id"`
	ToUser     User `gorm:"foreignKey:ToUserId;references:Id"`
}
