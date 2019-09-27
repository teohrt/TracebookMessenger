package dtos

// Node stores connection information and chat history
type Node struct {
	Name          string
	Port          string
	NodeAddress   string
	PeerAddresses []string
	ChatHistory   []string
	FirstPeer     string
}
