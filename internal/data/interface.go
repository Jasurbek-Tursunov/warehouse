package data

type Storage interface {
	MustConnect()
	HealthCheck()
	Close()
}
