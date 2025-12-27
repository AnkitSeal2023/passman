package utils

type Account struct {
	User string `json:"uname"`
	Pass string `json:"pass"`
}

type Response struct {
	Status bool      `json:"status"`
	Data   []Account `json:"data"`
}

func getEncryptedPass() {
}
