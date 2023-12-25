package generator

type Commands interface {
	GeneratePassword() string
	DumpToCsv() error
	DumpToDatabase(path string) error
	ExportToDropbox() error
}
