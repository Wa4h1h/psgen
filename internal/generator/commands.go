package generator

type Commands interface {
	GeneratePassword() error
	GetPassword(key string) (string, error)
	ExportToCsv() error
	ImportDatabase(path string) error
}
