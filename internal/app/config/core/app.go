package core

// App ..
type App struct {
	Secret string `default:"-" mapstructure:"secret"`
}
