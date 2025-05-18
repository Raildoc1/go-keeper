package protocol

type Creds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Entry struct {
	Metadata map[string]string `json:"metadata"`
	Data     []byte            `json:"data"`
}

type LoadRequest struct {
	ID int `json:"id"`
}
