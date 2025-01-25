package presenter

type Server interface {
	MustRun()
	GracefulStop()
}
