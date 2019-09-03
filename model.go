package main

type KeySecret struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type Login struct {
	Email    string `json:"email_address" form:"email_address"`
	Password string `json:"password" form:"password"`
}

type User struct {
	AccountID string `db:"id" goqu:"skipinsert,skipupdate"`
	Password  string `db:"password"`
	Email     string `db:"email"`
	Key       string `db:"key"`
	Secret    string `db:"secret"`
}

func (u *User) CheckPassword(pass string) bool {
	return u.Password == pass
}

type Credentials struct {
	JWT string `json:"jwt"`
}
