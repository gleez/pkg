// snowflake provides a very simple Twitter snowflake generator and parser.
package sid

// import (
// 	"errors"
// 	"fmt"
// 	"strconv"
// 	"sync"
// 	"time"
// "math/rand"
// )

// var (
// 	// Epoch is set to the twitter snowflake epoch of Nov 04 2010 01:42:54 UTC in milliseconds
// 	// You may customize this to set a different epoch for your application.
// 	Epoch int64 = 1288834974657

// 	// NodeBits holds the number of bits to use for Node
// 	// Remember, you have a total 22 bits to share between Node/Step
// 	NodeBits uint8 = 10

// 	// StepBits holds the number of bits to use for Step
// 	// Remember, you have a total 22 bits to share between Node/Step
// 	StepBits uint8 = 12
// )

// // A Node struct holds the basic information needed for a snowflake generator
// // node
// type Node struct {
// 	mu    sync.Mutex
// 	epoch time.Time
// 	time  int64
// 	node  int64
// 	step  int64

// 	nodeMax   int64
// 	nodeMask  int64
// 	stepMask  int64
// 	timeShift uint8
// 	nodeShift uint8
// }

// // NewNode returns a new snowflake node that can be used to generate snowflake
// // IDs
// func NewNode(node int64) (*Node, error) {
// 	n := Node{}
// 	n.node = node
// 	n.nodeMax = -1 ^ (-1 << NodeBits)
// 	n.nodeMask = n.nodeMax << StepBits
// 	n.stepMask = -1 ^ (-1 << StepBits)
// 	n.timeShift = NodeBits + StepBits
// 	n.nodeShift = StepBits

// 	if n.node < 0 || n.node > n.nodeMax {
// 		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
// 	}

// 	var curTime = time.Now()
// 	// add time.Duration to curTime to make sure we use the monotonic clock if available
// 	n.epoch = curTime.Add(time.Unix(Epoch/1000, (Epoch%1000)*1000000).Sub(curTime))

// 	return &n, nil
// }

// // Generate creates and returns a unique snowflake ID
// // To help guarantee uniqueness
// // - Make sure your system is keeping accurate system time
// // - Make sure you never have multiple nodes running with the same node ID
// func (n *Node) Generate() ID {

// 	n.mu.Lock()

// 	now := time.Since(n.epoch).Nanoseconds() / 1000000

// 	if now == n.time {
// 		n.step = (n.step + 1) & n.stepMask

// 		if n.step == 0 {
// 			for now <= n.time {
// 				now = time.Since(n.epoch).Nanoseconds() / 1000000
// 			}
// 		}
// 	} else {
// 		n.step = 0
// 	}

// 	n.time = now

// 	r := ID((now)<<n.timeShift |
// 		(n.node << n.nodeShift) |
// 		(n.step),
// 	)

// 	n.mu.Unlock()
// 	return r
// }

// func newIDGenerator() (*Node, error) {
// 	Epoch = toSonyflakeTime(time.Date(2020, 07, 01, 0, 0, 0, 0, time.UTC))

// 	rand.Seed(time.Now().UnixNano())
// 	nid := int64(rand.Int() & 1023)

// 	n, err := NewNode(nid)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return n, nil
// }

// func toSonyflakeTime(t time.Time) int64 {
// 	return t.UTC().UnixNano() / 1e7
// }
