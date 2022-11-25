package service

type LoginService interface {
	LoginUser(email string, password string) bool
}
type loginInformation struct {
	email    string
	password string
}

func NewLoginService() LoginService {
	return &loginInformation{
		email:    "a@a.com",
		password: "test123",
	}
}

func (info *loginInformation) LoginUser(email string, password string) bool {
	return info.email == email && info.password == password
}