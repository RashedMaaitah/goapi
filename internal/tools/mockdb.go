package tools

import "time"

type mockDB struct {
}

var mockLoginDetails = map[string]LoginDetails{
	"alex": {
		AuthToken: "123ABC",
		Username:  "alex",
	},
	"maria": {
		AuthToken: "456DEF",
		Username:  "maria",
	},
	"john": {
		AuthToken: "789GHI",
		Username:  "john",
	},
}

var mockCoinDetails = map[string]CoinDetails{
	"alex": {
		Coins:    1000,
		Username: "alex",
	},
	"maria": {
		Coins:    2500,
		Username: "maria",
	},
	"john": {
		Coins:    500,
		Username: "john",
	},
}

func (d *mockDB) GetUserLoginDetails(username string) *LoginDetails {

	time.Sleep(time.Second * 1)

	var clientData = LoginDetails{}

	clientData, ok := mockLoginDetails[username]

	if !ok {
		return nil
	}

	return &clientData
}

func (d *mockDB) GetUserCoins(username string) *CoinDetails {

	time.Sleep(time.Second * 1)
	var coinData = CoinDetails{}
	coinData, ok := mockCoinDetails[username]

	if !ok {
		return nil
	}
	return &coinData
}

func (d *mockDB) SetupDatabase() error {
	return nil
}
