package api

// Config is main config
type Config struct {
	Title  string
	Mode   string
	Addr   string
	Locale string

	Mysql struct {
		Addr     string
		User     string
		Password string
		Database string
	}
}
