package models

type Secret struct {
	Host     string `json:"host"`
	UserName string `json:"username"`
	Password string `json:"password"`
	JWTSign  string `json:"jwtsign"`
	DataBase string `json:"database"`
}
