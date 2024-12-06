package core

// App ..
type App struct {
	Secret    string `default:"-"    mapstructure:"secret"`
	IsShowLog string `default:"true" mapstructure:"is_show_log"`
}
