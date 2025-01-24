package usecase

type AuthService interface {
	Register(username, password string) (string, error)
	Login(username, password string) (string, error)
}
