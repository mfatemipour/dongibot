package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type CoeffPair struct {
	Credit float64 `json:"credit"`
	Debt   float64 `json:"debt"`
}

type ShareCoeffs map[uint]CoeffPair

func (sc *ShareCoeffs) GetCoeffSum() (credit, debt float64) {
	for _, v := range *sc {
		credit += v.Credit
		debt += v.Debt
	}
	return credit, debt
}

func (sc *ShareCoeffs) Scan(value interface{}) error {
	fmt.Printf("%s\n", value.(string))
	jsonStr, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	if err := json.Unmarshal([]byte(jsonStr), sc); err != nil {
		return err
	}
	return nil
}

func (sc ShareCoeffs) Value() (driver.Value, error) {
	e, err := json.Marshal(sc)
	if err != nil {
		return nil, err
	}
	return string(e), nil
}

type User struct {
	gorm.Model
	ID        uint
	Name      string
	IsAdmin   bool
	DongUsers []DongUser
}

type DongUser struct {
	gorm.Model
	ID           uint
	UserID       uint `gorm:"UNIQUE_INDEX:compositeindex"`
	DongID       uint `gorm:"UNIQUE_INDEX:compositeindex"`
	IsDongAdmin  bool
	Transactions []Transaction
}

type Transaction struct {
	gorm.Model
	ID          uint
	DongUserID  uint
	DongID      uint
	Expense     int
	ShareCoeffs ShareCoeffs `gorm:"type:text"`
	Desc        string
}

type Dong struct {
	gorm.Model
	ID           uint
	Name         string
	Desc         string
	Transactions []Transaction
}
