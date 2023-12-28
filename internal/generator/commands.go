package generator

type Commands interface {
	GeneratePassword() error
	GetPassword(key string) (string, error)
	Export() error
	Import(csvPath string, concurrency int, withEncryption bool) error
}
