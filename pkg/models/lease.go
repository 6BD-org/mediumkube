package models

// PeerLease represents peer nodes in a mediumkube cluster
type PeerLease struct {
	Cidr      string `json:"cidr"`      // Cidr of the node
	Timestamp int64  `json:"timestamp"` // Current unix timestamp
	TTL       int64  `json:"ttl"`       // Time to live
}
