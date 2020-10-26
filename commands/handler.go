package commands

type Handler interface {
	Handle(args []string)
	Help()
	Desc() string
}
