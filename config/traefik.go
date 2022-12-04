package config

type Traefik struct {
	Enable    bool
	Schema    string
	Host      string
	BasicAuth bool
	Username  string
	Password  string
}
