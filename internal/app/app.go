package app

// App wires dependencies and exposes workflows to the CLI.
type App struct{}

// New builds a new application container.
func New() *App {
	return &App{}
}
