package controllers

import (
	"awesomeProject/dao"
	"awesomeProject/services"
	u "awesomeProject/utils"
	_ "encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

// BuyItem - эндпоинт для покупки мерча пользователем.
var BuyItem = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var merchType = vars["item"]
	merch, err := services.GetMerch(merchType)
	if err != nil {
		u.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

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

	// Проверяем баланс пользователя.
	if user.Balance < merch.Price {
		u.RespondWithError(w, http.StatusBadRequest, "not enough coins")
	}

	// Оформляем покупку как транзакцию в БД.
	err = services.GetDB().Transaction(func(tx *gorm.DB) error {
		purchase := dao.Purchase{UserId: user.Id, MerchId: merch.Id}
		if err := tx.Create(&purchase).Error; err != nil {
			return err
		}

		user.Balance -= merch.Price
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		u.RespondWithError(w, http.StatusInternalServerError, "error while buying an item")
		return
	}

	u.RespondJSON(w, http.StatusOK, nil)
}
