package dto

type Creds struct {
	Username string
	Password string
}

type Entry struct {
	Metadata map[string]string
	Data     []byte
}
