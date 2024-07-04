package metadata

type PayloadMetadata struct {
	Paging  PagingPayload `json:"paging"`
	Metrics []string      `json:"metrics"`
}

type PagingPayload struct {
	From int `json:"from"`
	To   int `json:"to"`
}
