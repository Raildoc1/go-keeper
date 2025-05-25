package options

type Option interface {
	GetHeaders() (map[string]string, error)
}
