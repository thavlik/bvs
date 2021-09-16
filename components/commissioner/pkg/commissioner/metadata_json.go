package commissioner

import (
	"fmt"
)

func metadataJson(id, policyID string, timestamp int64) string {
	return fmt.Sprintf(`{
	"721": {
		"%s": {
			"Vote": {
				"description": "This is a vote in an election for the mayor of Bikini Bottom.",
				"name": "Blockchain Voting Systems NFT vote",
				"id": "%s",
				"timestamp": %d
			}
		}
	}
}`, policyID, id, timestamp)
}
