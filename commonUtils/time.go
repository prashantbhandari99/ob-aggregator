package commonUtils

import (
	"time"
)

func CurrentTimestampMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
