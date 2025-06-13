package protocol

type Creds struct {
	Username string `json:"username"`
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
