package entity

type User struct {
	ID        string
	Name      string
	Email     string
	Token     string
	Password  string
	CreatedAt int64
}

// 去除敏感信息
func (u *User) Info() *User {
	u.Password = ""
	u.Token = ""
	return u
}
