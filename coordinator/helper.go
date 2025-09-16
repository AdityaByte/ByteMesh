package coordinator

import (
	"github.com/AdityaByte/bytemesh/payload"
)

// CheckHealthofDataNode - If there were no node alive then it will return false otherwise returns true.
func isNodeHealthGood(registeredNodes payload.RegisteredDataNodes) bool {
	if len(registeredNodes.Nodes) == 0 {
		return false
	}
	return true
}
