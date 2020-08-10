package sid

// https://github.com/RO-29/snowflake-golang
import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	timeCustom = "2020-08-01T00:00:00+00:00"

	totalBITS    = 64
	epochBITS    = 42 // time - Epoch timestamp in milliseconds precision - 42 bits
	nodeIDBITS   = 10 // configured machine id - 10 bits. This gives us 1024 nodes/machines
	sequenceBITS = 12 // sequence number - 12 bits - rolls over every 4096 per machine (with protection to avoid rollover in the same ms)
)

var (
	// Custom Epoch (in milliseconds) (August 1, 2020 Midnight UTC = 2020-08-01T00:00:00Z)
	customEPOCH int64
	nodeID      int

	maxNodeID   = (int)(math.Pow(2, float64(nodeIDBITS)) - 1)
	maxSequence = (int)(math.Pow(2, float64(sequenceBITS)) - 1)

	sid *Snowflake
)

func init() {
	rand.Seed(time.Now().UnixNano())
	timeMustParse()
	nodeIDGenerator()

	if sid == nil {
		sid = NewSnowFlake()
	}
}

func New() ID {
	return sid.Generate()
}

//NewSnowFlake service init
func NewSnowFlake() *Snowflake {
	return &Snowflake{
		lastTimestamp: -1,
	}
}

//Snowflake service ...
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	sequence      int64
}

func (s *Snowflake) Generate() ID {
	if id, err := s.GenerateUniqueSequenceID(); err == nil {
		return ID(id)
	}

	return 0
}

// GenerateUniqueSequenceID generates unique id ...
func (s *Snowflake) GenerateUniqueSequenceID() (int64, error) {

	currentTimeStamp, err := s.generateCurrentTimeSequence()
	if err != nil {
		return 0, err
	}

	// first 42 bits of our ID will be filled with the epoch timestamp. left-shift to achieve this
	id := currentTimeStamp << uint(totalBITS-epochBITS)

	// fill the next 10 bits with the node ID.
	id |= int64(nodeID << uint(totalBITS-epochBITS-nodeIDBITS))

	// last 12 bits with the local counter.
	id |= s.sequence
	return id, nil

}

func (s *Snowflake) generateCurrentTimeSequence() (int64, error) {

	s.mu.Lock()
	currentTimeStamp, err := s.getCurrentTimeStamp()
	if err != nil {
		return 0, err
	}
	s.lastTimestamp = currentTimeStamp
	s.mu.Unlock()

	return currentTimeStamp, nil
}

func (s *Snowflake) getCurrentTimeStamp() (int64, error) {
	currentTimeStamp := getTimeStampMilli()

	if currentTimeStamp < s.lastTimestamp {
		return 0, errors.New("sid: invalid system clock")
	}

	if currentTimeStamp > s.lastTimestamp {
		// reset sequence to start with zero for the next millisecond
		s.sequence = 0
		return currentTimeStamp, nil
	}

	s.sequence = (s.sequence + 1) & int64(maxSequence)
	if s.sequence != 0 {
		return currentTimeStamp, nil
	}

	// Sequence Exhausted, wait till next millisecond.
	return s.waitNextMillis(currentTimeStamp), nil

}

// Block and wait till next millisecond
func (s *Snowflake) waitNextMillis(currentTimeStamp int64) int64 {
	for currentTimeStamp == s.lastTimestamp {
		currentTimeStamp = getTimeStampMilli()
	}
	return currentTimeStamp
}

// Get the current timestamp in milliseconds, adjust for the custom epoch.
func getTimeStampMilli() int64 {
	return time.Now().UnixNano()/1e6 - customEPOCH
}

func timeMustParse() {
	timeObj, err := time.Parse(time.RFC3339, timeCustom)
	if err != nil {
		panic(err)
	}

	customEPOCH = timeObj.UnixNano() / 1e6
}

func nodeIDGenerator() {
	//Make sure to generate it once only ...
	if nodeID != 0 {
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("sid: error generating nodeID::", err.Error())

		//In case it fails to get mac address, generate the random number,
		//gloabbly seeded with time.Now()
		nodeID = rand.Int() & maxNodeID
	}

	nodeID = hashCode(hostname) & maxNodeID

	fmt.Printf("sid: NodeId %d - MaxNodeId %d \n", nodeID, maxNodeID)
}

func hashCode(s string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return int(h.Sum32())
}

// NewNodeByHostname is a convenience method which creates a new Node based
// off a hash of the machine's hostname.
func NewNodeByHostname() int {
	name, err := os.Hostname()
	if err != nil {
		return 0
	}

	hash := md5.Sum([]byte(name))
	id := binary.BigEndian.Uint64(hash[:]) & 0x3FF // mask to first 10 bits, max of 1023

	return int(id)
}
