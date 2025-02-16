package dao

type Merch struct {
	Id    uint   `json:"id" gorm:"primaryKey"`
	Type  string `json:"type"`
	Price int    `json:"price"`
}

// MerchList - список существующего мерча.
var MerchList = []Merch{
	{1, "t-shirt", 80},
	{2, "cup", 20},
	{3, "book", 50},
	{4, "pen", 10},
	{5, "powerbank", 200},
	{6, "hoody", 300},
	{7, "umbrella", 200},
	{8, "socks", 10},
	{9, "wallet", 50},
	{10, "pink-hoody", 500},
}
