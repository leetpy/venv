package internal

type Node struct {
	ID      string   `json:"id"`
	Type    NodeType `json:"type"`
	IPMI    string   `json:"ipmi"`
	Remarks string   `json:"remarks"`
}

type Port struct {
	ID       string `json:"id"`
	Mac      string `json:"mac"`
	PortId   string `json:"port_id"`
	SwitchId string `json:"switch_id"`
}

type NodeType struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
