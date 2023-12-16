package generator

type Commands interface {
	GeneratePassword() string
}
