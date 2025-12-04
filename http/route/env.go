package route

type Environment struct {
	Title string
}

func NewEnv() *Environment {
	env := new(Environment)

	env.Title = "Hyperuplink"

	return env
}
