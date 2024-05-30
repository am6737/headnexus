package entity

type User struct {
	ID           string
	Name         string
	Email        string
	Token        string
	Verification string
	Status       UserStatus
	Password     string
	CreatedAt    int64
	LastLoginAt  int64
}

// 去除敏感信息
func (u *User) Info() *User {
	//u.Password = ""
	u.Token = ""
	return u
}

type UserStatus uint

const (
	Normal  UserStatus = iota + 1 //正常
	Disable                       //禁用
)

type FindOptions struct {
	Email        string
	Verification string
	Token        string
}
