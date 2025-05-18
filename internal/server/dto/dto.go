package dto

type Creds struct {
	Username string
	Password string
}

type Entry struct {
	ID       string
	Metadata map[string]string
	Data     []byte
}
