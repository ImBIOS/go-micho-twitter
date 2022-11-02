package responses

type Error struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
