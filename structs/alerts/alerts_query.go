package alerts

type AlertQuery struct {
	Alerts []Alert
}

type Alert struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
}
