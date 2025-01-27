package di

type StorageType int

const (
	StorageDefault  StorageType = StorageMock
	StorageMock                 = 0
	StoragePostgres             = 1
	StorageMySQL                = 2
)

type ServerType int

const (
	ServerDefault ServerType = ServerGin
	ServerGin                = 0
	ServerStd                = 1
	ServerGRPC               = 2
)
