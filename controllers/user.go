package controllers

import (
	"awesomeProject/dao"
	"awesomeProject/dto"
	"awesomeProject/services"
	u "awesomeProject/utils"
	"encoding/json"
	_ "github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

// SendCoins - эндпоинт для передачи монет.
var SendCoins = func(w http.ResponseWriter, r *http.Request) {

	var sendCoinRequest dto.SendCoinRequest
	err := json.NewDecoder(r.Body).Decode(&sendCoinRequest)
	if err != nil {
		u.RespondWithError(w, http.StatusBadRequest, "Invalid request body.\n"+err.Error())
		return
	}

	username := r.Context().Value("user")
	if username == nil {
		u.RespondWithError(w, http.StatusBadRequest, "no user in context")
		return
	}

	user, err := services.GetUser(username.(string))
	// Проверяем пользователя-отправителя.
	if err != nil {
		u.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	toUser, err := services.GetUser(sendCoinRequest.ToUser)
	// Проверяем пользователя-получателя.
	if err != nil {
		u.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Проверяем баланс пользователя.
	if user.Balance < sendCoinRequest.Amount {
		u.RespondWithError(w, http.StatusBadRequest, "not enough coins")
	}

	// Проверяем получателя.
	if user.Username == toUser.Username {
		u.RespondWithError(w, http.StatusBadRequest, "cannot send coins to yourself")
	}

	// Оформляем перевод как транзакцию в БД.
	err = services.GetDB().Transaction(func(tx *gorm.DB) error {
		user.Balance -= sendCoinRequest.Amount
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		toUser.Balance += sendCoinRequest.Amount
		if err := tx.Save(&toUser).Error; err != nil {
			return err
		}

		transaction := dao.Transaction{FromUserId: user.Id, ToUserId: toUser.Id, Amount: sendCoinRequest.Amount}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		u.RespondWithError(w, http.StatusInternalServerError, "error while sending coins")
		return
	}

	u.RespondJSON(w, http.StatusOK, nil)
}

// GetInfo - эндпоинт для получения информации о движении монет.
var GetInfo = func(w http.ResponseWriter, r *http.Request) {

	username := r.Context().Value("user")
	if username == nil {
		u.RespondWithError(w, http.StatusBadRequest, "no user in context")
		return
	}

	user, err := services.GetUser(username.(string))
	// Проверяем пользователя.
	if err != nil {
		u.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var (
		coins                = user.Balance
		items                []dto.Item
		receivedTransactions []dto.ReceivedTransaction
		sentTransactions     []dto.SentTransaction
	)

	// Получаем покупки пользователя.
	err = services.GetDB().Table("purchases").
		Select("merches.type, COUNT(*) as quantity").
		Joins("INNER JOIN merches ON purchases.merch_id = merches.id").
		Where("purchases.user_id = ?", user.Id).
		Group("merches.type").
		Scan(&items).Error
	if err != nil {
		u.RespondWithError(w, http.StatusInternalServerError, "error while getting info")
		return
	}

	// Получаем переводы пользователя.
	err = services.GetDB().Table("transactions").
		Select("users.username AS to_user, SUM(transactions.amount) AS amount").
		Joins("INNER JOIN users ON transactions.to_user_id = users.id").
		Where("transactions.from_user_id = ?", user.Id).
		Group("users.username").
		Scan(&sentTransactions).Error
	if err != nil {
		u.RespondWithError(w, http.StatusInternalServerError, "error while getting info")
		return
	}

	// Получаем переводы пользователю.
	err = services.GetDB().Table("transactions").
		Select("users.username AS from_user, SUM(transactions.amount) AS amount").
		Joins("INNER JOIN users ON transactions.from_user_id = users.id").
		Where("transactions.to_user_id = ?", user.Id).
		Group("users.username").
		Scan(&receivedTransactions).Error
	if err != nil {
		u.RespondWithError(w, http.StatusInternalServerError, "error while getting info")
		return
	}

	coinHistory := dto.CoinHistory{Received: receivedTransactions, Sent: sentTransactions}
	infoResponse := dto.InfoResponse{Coins: coins, Inventory: items, CoinHistory: coinHistory}

	u.RespondJSON(w, http.StatusOK, infoResponse)
}
