package protocol

type Creds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Entry struct {
	ID       string            `json:"id"`
	Metadata map[string]string `json:"metadata"`
	Data     []byte            `json:"data"`
}

type LoadRequest struct {
	ID string `json:"id"`
}
