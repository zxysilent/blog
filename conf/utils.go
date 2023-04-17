package conf

import (
	"fmt"
)

func IsProd() bool {
	return App.Mode != "dev"
}
func IsDebug() bool {
	return App.Mode == "dev"
}

func Addr() string {
	return fmt.Sprintf("%s:%d", App.Host, App.Port)
}

func Endpoint() string {
	return App.Endpoint
}
