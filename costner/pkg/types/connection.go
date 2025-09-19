package types

type Connection struct {
	ID         string `json:"id"`
	SourceNode string `json:"source_node"`
	SourcePort string `json:"source_port"`
	TargetNode string `json:"target_node"`
	TargetPort string `json:"target_port"`
}

type ConnectionPoint struct {
	NodeID string `json:"node_id"`
	Port   string `json:"port"`
}