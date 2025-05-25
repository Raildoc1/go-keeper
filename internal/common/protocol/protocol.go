package protocol

type Creds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type StoreRequest struct {
	GUID  string `json:"guid"`
	Entry Entry  `json:"entry"`
}

type Entry struct {
	Metadata map[string]string `json:"metadata"`
	Data     []byte            `json:"data"`
}

type LoadRequest struct {
	GUID string `json:"guid"`
}
