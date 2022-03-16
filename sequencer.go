package main

import (
	"fmt"
	"hash/crc32"
	"math"
	"math/rand"
	"net"
	"sync"
	"time"
)

// Adapted from https://www.callicoder.com/distributed-unique-id-sequence-number-generator/
// (which was adapter from Twitter Snowflake).

const (
	UnusedBits   = 1 // Sign bit, Unused (always set to 0)
	EpochBits    = 41
	NodeIdBits   = 10
	SequenceBits = 12
	CustomEpoch  = 1420070400000 // Custom Epoch (January 1, 2015 Midnight UTC = 2015-01-01T00:00:00Z)
)

type sequencer struct {
	nodeId        int
	lastTimestamp int64
	sequence      int64
	maxNodeId     int64
	maxSequence   int64
	mutex         *sync.Mutex
}

func baseSequencer() sequencer {
	s := sequencer{}
	s.maxNodeId = int64(math.Pow(2, NodeIdBits) - 1)
	s.maxSequence = int64(math.Pow(2, SequenceBits) - 1)
	return s
}

func NewSequencerWithNodeId(nodeId int) (sequencer, error) {
	s := baseSequencer()

	if nodeId < 0 || int64(s.nodeId) > s.maxNodeId {
		return s, fmt.Errorf("NodeId must be between %d and %d", 0, s.maxNodeId)
	}
	s.nodeId = nodeId
	s.lastTimestamp = -1
	s.sequence = 0
	s.mutex = &sync.Mutex{}
	return s, nil
}

func NewSequencer() (sequencer, error) {
	s := baseSequencer()

	nodeId, err := s.createNodeId()
	if err != nil {
		return s, err
	}
	return NewSequencerWithNodeId(nodeId)
}

func (s *sequencer) NextId() (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	currentTimestamp := s.timestamp()

	if currentTimestamp < s.lastTimestamp {
		return 0, fmt.Errorf("invalid system clock!")
	}

	if currentTimestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & s.maxSequence
		if s.sequence == 0 {
			// Sequence Exhausted, wait till next millisecond.
			currentTimestamp = s.waitNextMillis(currentTimestamp)
		}
	} else {
		// reset sequence to start with zero for the next millisecond
		s.sequence = 0
	}

	s.lastTimestamp = currentTimestamp

	id := currentTimestamp << (NodeIdBits + SequenceBits)
	id |= int64(s.nodeId) << int64(SequenceBits)
	id |= s.sequence
	return id, nil
}

// Get current timestamp in milliseconds, adjust for the custom epoch.
func (s *sequencer) timestamp() int64 {
	return time.Now().UnixMilli() - CustomEpoch
}

func (s *sequencer) waitNextMillis(currentTimestamp int64) int64 {
	for {
		if currentTimestamp != s.lastTimestamp {
			break
		}
		currentTimestamp = s.timestamp()
		time.Sleep(time.Millisecond)
	}
	return currentTimestamp
}

func (s *sequencer) createNodeId() (int, error) {
	nodeId := 0
	ifaces, err := net.Interfaces()
	if err != nil {
		return 0, err
	}
	hash := ""
	for _, iface := range ifaces {
		//log.Printf("IFACE : %#v", iface)
		mac := iface.HardwareAddr.String()
		//log.Printf("IFACE : [mac=%s] %#v", mac, iface)
		if mac != "" {
			if hash != "" {
				hash += ":"
			}
			hash += mac
		}
		//log.Printf("hash = %s, %d", string(hash), crc32.ChecksumIEEE([]byte(hash)))
	}
	if hash != "" {
		nodeId = int(crc32.ChecksumIEEE([]byte(hash)))
	} else {
		nodeId = int(rand.Int63())
	}
	if nodeId < 0 {
		nodeId = -nodeId
	}
	//log.Printf("nodeId = %d, s.maxNodeId = %d", nodeId, s.maxNodeId)
	nodeId = nodeId & int(s.maxNodeId)
	return nodeId, nil
}
