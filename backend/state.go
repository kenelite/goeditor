package backend

type State struct {
	CurrentFile string
}

func NewState() *State {
	return &State{}
}
