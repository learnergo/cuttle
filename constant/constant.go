package constant

type NodeType string

const (
	Peer    NodeType = "peer"
	Orderer NodeType = "orderer"
	Admin   NodeType = "admin"
	User    NodeType = "user"
)
