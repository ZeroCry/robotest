package gravity

import (
	"time"
)

// default timeout to wait for cloud-init to complete
var cloudInitTimeout = time.Minute * 30

// default timeout to wait for clocks to synchronize between nodes
var clockSyncTimeout = time.Minute * 15