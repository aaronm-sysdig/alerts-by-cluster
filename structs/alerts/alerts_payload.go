package alerts

type PayloadAlert struct {
	Enabled                bool            `json:"enabled"`
	Type                   string          `json:"type"`
	Name                   string          `json:"name"`
	Description            string          `json:"description"`
	Scope                  string          `json:"scope"`
	Repositories           []string        `json:"repositories"`
	Triggers               PayloadTriggers `json:"triggers"`
	Autoscan               bool            `json:"autoscan"`
	OnlyPassFail           bool            `json:"onlyPassFail"`
	NotificationChannelIds []string        `json:"notificationChannelIds"`
}

type PayloadTriggers struct {
	Unscanned      bool `json:"unscanned"`
	AnalysisUpdate bool `json:"analysis_update"`
	VulnUpdate     bool `json:"vuln_update"`
	PolicyEval     bool `json:"policy_eval"`
}
