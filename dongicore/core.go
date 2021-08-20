package dongicore

import (
	"fmt"

	"github.com/mfatemipour/dongibot/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DongiCore struct {
	// dongs     map[uint]*database.Dong
	dbHandler *database.DBHandler
}

func NewDongiCore(dbHandler *database.DBHandler) (*DongiCore, error) {
	core := new(DongiCore)
	core.dbHandler = dbHandler

	// dongs := make([]database.Dong, 0)
	// res := dbHandler.DB.Find(&dongs)
	// logrus.Debugf("loading db: %v", res)

	// for _, dong := range dongs {
	// 	core.dongs[dong.ID] = &dong
	// }
	return core, nil
}

func (core *DongiCore) AddUser(user *database.User) error {
	res := core.dbHandler.DB.Create(user)
	if res.Error != nil {
		return res.Error
	}
	logrus.Debugf("Add user db api succeeded, rows affected: %d", res.RowsAffected)
	return nil
}

func (core *DongiCore) UpdateUser() {

}

func (core *DongiCore) DeleteUser() {

}

func (core *DongiCore) AddDongUser(dongUser *database.DongUser) error {
	var user database.User
	var dong database.Dong
	var tx *gorm.DB
	tx = core.dbHandler.DB.First(&user, dongUser.UserID)
	if tx.Error != nil {
		return tx.Error
	}
	tx = core.dbHandler.DB.First(&dong, dongUser.DongID)
	if tx.Error != nil {
		return tx.Error
	}
	tx = core.dbHandler.DB.Create(&dongUser)
	if tx.Error != nil {
		return tx.Error
	}
	logrus.Debugf("adding dongUser succeeded, rows affected: %d", tx.RowsAffected)
	return nil
}

// func (core *DongiCore) UpdateDongUser() {

// }

// func (core *DongiCore) DeleteDongUser() {

// }

func (core *DongiCore) AddDong(dong *database.Dong) error {
	res := core.dbHandler.DB.Create(dong)
	if res.Error != nil {
		return res.Error
	}
	logrus.Debugf("Add dong db api succeeded, rows affected: %d", res.RowsAffected)
	return nil
}

func (core *DongiCore) UpdateDong() {

}

func (core *DongiCore) DeleteDong() {

}

func (core *DongiCore) AddTransaction(transaction *database.Transaction) error {
	var tx *gorm.DB
	dong := database.Dong{}
	tx = core.dbHandler.DB.First(&dong, transaction.DongID)
	if tx.Error != nil {
		return tx.Error
	}

	dongusers := make([]database.DongUser, 0)
	tx = core.dbHandler.DB.Where("dong_id == ?", transaction.DongID).Find(&dongusers)
	if tx.Error != nil {
		return tx.Error
	}

	cSum, dSum := transaction.ShareCoeffs.GetCoeffSum()
	if cSum == 0 || dSum == 0 {
		return fmt.Errorf("transaction is invalid, credis sum: %f, dept sum: %f", cSum, dSum)
	}

	for k := range transaction.ShareCoeffs {
		valid := false
		for _, dongUser := range dongusers {
			if dongUser.UserID == k {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("dong %d does not contain user [%d]", dong.ID, k)
		}
	}
	tx = core.dbHandler.DB.Create(&transaction)
	if tx.Error != nil {
		return tx.Error
	}
	logrus.Debugf("adding transaction by dong-user %d to dong %d succeeded, rows affected: %d",
		transaction.DongUserID, transaction.DongID, tx.RowsAffected)
	return nil
}

func (core *DongiCore) UpdateTransaction() {

}

func (core *DongiCore) DeleteTransaction() {

}

func (core *DongiCore) GenerateInvoce(dongID uint) (map[uint]float64, error) {
	var tx *gorm.DB
	var dong database.Dong
	tx = core.dbHandler.DB.Preload("Transactions").Find(&dong, dongID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	dongUsers := make([]database.DongUser, 0)
	tx = core.dbHandler.DB.Where("dong_id == ?", dongID).Find(&dongUsers)
	if tx.Error != nil {
		return nil, tx.Error
	}

	invoice := make(map[uint]float64)

	for _, dongUser := range dongUsers {
		invoice[dongUser.UserID] = 0
	}

	expenseSum := 0

	for _, transaction := range dong.Transactions {
		expenseSum += transaction.Expense
		cSum, dSum := transaction.ShareCoeffs.GetCoeffSum()
		for k, coeffs := range transaction.ShareCoeffs {
			share := (coeffs.Credit/cSum)*float64(transaction.Expense) -
				(coeffs.Debt/dSum)*float64(transaction.Expense)
			preShare := invoice[k]
			invoice[k] = preShare + share
		}
	}
	return invoice, nil
}
