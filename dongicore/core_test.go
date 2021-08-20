package dongicore

import (
	"fmt"
	"os"
	"testing"

	"github.com/mfatemipour/dongibot/database"
	"github.com/stretchr/testify/assert"
)

func TestCore(t *testing.T) {
	db, err := database.NewDB("/tmp/test.db")
	assert.Equal(t, nil, err)
	fmt.Printf("%v\n", db)
	core, err := NewDongiCore(db)
	assert.Equal(t, nil, err)
	dong := database.Dong{Name: "dong", Desc: "desc"}
	user1 := database.User{Name: "user1", IsAdmin: false}
	user2 := database.User{Name: "user2", IsAdmin: false}
	err = core.AddDong(&dong)
	assert.Equal(t, nil, err)
	err = core.AddUser(&user1)
	assert.Equal(t, nil, err)
	err = core.AddUser(&user2)
	assert.Equal(t, nil, err)

	dongUser1 := database.DongUser{UserID: user1.ID, DongID: dong.ID}
	dongUser2 := database.DongUser{UserID: user2.ID, DongID: dong.ID}
	err = core.AddDongUser(&dongUser1)
	assert.Equal(t, nil, err)
	err = core.AddDongUser(&dongUser2)
	assert.Equal(t, nil, err)

	trx1 := database.Transaction{DongUserID: dongUser1.ID, DongID: dong.ID, Expense: 1000, ShareCoeffs: database.ShareCoeffs{
		user1.ID: database.CoeffPair{Credit: 1, Debt: 1},
		user2.ID: database.CoeffPair{Credit: 0, Debt: 1},
	}}
	err = core.AddTransaction(&trx1)
	assert.Equal(t, nil, err)

	err = os.Remove("/tmp/test.db")
	assert.Equal(t, nil, err)
}

func TestSimpleInvoice(t *testing.T) {
	db, err := database.NewDB("/tmp/test.db")
	assert.Equal(t, nil, err)
	fmt.Printf("%v\n", db)
	core, err := NewDongiCore(db)
	assert.Equal(t, nil, err)
	dong := database.Dong{Name: "dong", Desc: "desc"}
	user1 := database.User{Name: "user1", IsAdmin: false}
	user2 := database.User{Name: "user2", IsAdmin: false}
	err = core.AddDong(&dong)
	assert.Equal(t, nil, err)
	err = core.AddUser(&user1)
	assert.Equal(t, nil, err)
	err = core.AddUser(&user2)
	assert.Equal(t, nil, err)

	dongUser1 := database.DongUser{UserID: user1.ID, DongID: dong.ID}
	dongUser2 := database.DongUser{UserID: user2.ID, DongID: dong.ID}
	err = core.AddDongUser(&dongUser1)
	assert.Equal(t, nil, err)
	err = core.AddDongUser(&dongUser2)
	assert.Equal(t, nil, err)

	trx1 := database.Transaction{DongUserID: dongUser1.ID, DongID: dong.ID,
		Expense: 1000, ShareCoeffs: database.ShareCoeffs{
			user1.ID: database.CoeffPair{Credit: 1, Debt: 1},
			user2.ID: database.CoeffPair{Credit: 0, Debt: 1},
		}}
	err = core.AddTransaction(&trx1)
	assert.Equal(t, nil, err)

	trx2 := database.Transaction{DongUserID: dongUser1.ID, DongID: dong.ID,
		Expense: 5000, ShareCoeffs: database.ShareCoeffs{
			user1.ID: database.CoeffPair{Credit: 0, Debt: 1},
			user2.ID: database.CoeffPair{Credit: 1, Debt: 1},
		}}
	err = core.AddTransaction(&trx2)
	assert.Equal(t, nil, err)

	invoice, err := core.GenerateInvoce(dong.ID)
	assert.Equal(t, nil, err)

	assert.Equal(t, float64(-2000), invoice[user1.ID])
	assert.Equal(t, float64(2000), invoice[user2.ID])

	err = os.Remove("/tmp/test.db")
	assert.Equal(t, nil, err)
}
