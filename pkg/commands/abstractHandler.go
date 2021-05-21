package commands

// Handler abstract handler.
type Handler interface {
	Handle(args []string)
	Help()
	Desc() string
}
