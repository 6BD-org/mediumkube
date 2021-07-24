package models

// PeerLease represents peer nodes in a mediumkube cluster
type PeerLease struct {
	Id        string `json:"id"`        // Id of node
	Cidr      string `json:"cidr"`      // Cidr of the node
	Timestamp int64  `json:"timestamp"` // Current unix timestamp of last refreshed time
	TTL       int64  `json:"ttl"`       // Time to live
}
