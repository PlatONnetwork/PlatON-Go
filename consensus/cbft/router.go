// Package cbft implements  a concrete consensus engines.
package cbft

const (

//DEFAULT_FANOUT = 5	// gossip protocol default value of fan-out.
)

// Router implements message forwarding.
type router struct {
}

// NewRouter creates a new router. It is mainly
// used for message forwarding
func NewRouter() *router {
	return nil
}
