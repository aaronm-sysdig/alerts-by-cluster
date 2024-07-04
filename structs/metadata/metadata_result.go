package metadata

type ResultMetadata struct {
	Metrics []string                `json:"metrics"`
	Time    TimeRangeMetadataResult `json:"time"`
	Data    []DataMetadataResult    `json:"data"`
	Paging  PagingMetadataResult    `json:"paging"`
}

type TimeRangeMetadataResult struct {
	From     int64 `json:"from"`
	To       int64 `json:"to"`
	Sampling int64 `json:"sampling"`
}

type DataMetadataResult struct {
	KubernetesClusterName string `json:"kubernetes.cluster.name"`
}

type PagingMetadataResult struct {
	From  int `json:"from"`
	To    int `json:"to"`
	Total int `json:"total"`
}
