package generator

type Commands interface {
	GeneratePassword() error
	GetPassword(key string) (string, error)
	Export(path string) error
	Import(csvPath string, concurrency int, withEncryption bool) error
}
