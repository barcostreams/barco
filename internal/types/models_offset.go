package types

import (
	"fmt"
	"math"
	"time"
)

type OffsetCommitType int
type CompareResult int

const (
	OffsetCommitNone OffsetCommitType = iota
	OffsetCommitLocal
	OffsetCommitAll
)

const (
	CompareEqual CompareResult = iota
	CompareLessThan
	CompareGreaterThan
)

const OffsetCompleted = math.MaxInt64

type OffsetSource struct {
	Id        GenId `json:"id"` // Gen id of the source
	Timestamp int64 `json:"ts"` // Timestamp in Unix Micros
}

func NewOffsetSource(id GenId) OffsetSource {
	return OffsetSource{
		Id:        id,
		Timestamp: time.Now().UnixMicro(),
	}
}

// Represents a topic offset for a given token.
type Offset struct {
	Version     GenVersion   `json:"version"`     // Generation version of the offset
	ClusterSize int          `json:"clusterSize"` // Cluster size of the gen version
	Offset      int64        `json:"value"`       // Numerical offset value
	Token       Token        `json:"token"`       // The start token of the offset generation
	Index       RangeIndex   `json:"index"`       // The range index of the offset
	Source      OffsetSource `json:"source"`      // The point-in-time when the offset was recorded.
}

// Checks whether a consumer that is assigned the provided range can read a past offset.
// Used when ranges don't match
func (o *Offset) IsAssignedToConsumer(currentToken Token, currentIndex RangeIndex, clusterSize int) bool {
	// checks whether the current token&index contains the offset's token&index
	// TODO: IMPLEMENT
	return false
}

func (o *Offset) IsBrokerAssigned(leaderTokens []Token, clusterSize int) bool {
	// TODO: IMPLEMENT
	return false
}

func (o *Offset) Compare(other *Offset) CompareResult {
	// TODO: REEVALUATE WHETHER IT SHOULD BE USED
	if other == nil {
		return CompareGreaterThan
	}

	if o.Version < other.Version {
		return CompareLessThan
	}
	if o.Version > other.Version {
		return CompareGreaterThan
	}
	if o.Offset < other.Offset {
		return CompareLessThan
	}
	if o.Offset > other.Offset {
		return CompareGreaterThan
	}
	return CompareEqual
}

func (o *Offset) String() string {
	return fmt.Sprintf("v%d %d", o.Version, o.Offset)
}

// Represents an identifier of an offset to be persisted
type OffsetStoreKey struct {
	Group string `json:"group"`
	Topic string `json:"topic"`
}

// Represents an identifier and a value of an offset
type OffsetStoreKeyValue struct {
	Key   OffsetStoreKey `json:"key"`
	Value Offset         `json:"value"`
}

// Represents a local view of the consumer group offsets
type OffsetState interface {
	Initializer
	fmt.Stringer

	// Gets the offset value for a given group and token.
	// Returns nil when not found
	//
	// The caller MUST check whether the current broker can serve the data when ranges don't match
	// The caller MUST check whether the consumer is assigned when ranges don't match
	Get(group string, topic string, token Token, rangeIndex RangeIndex, clusterSize int) (offset *Offset, rangesMatch bool)

	// Sets the known offset value in memory, optionally committing it to the data store
	Set(group string, topic string, value Offset, commit OffsetCommitType) bool

	// Reads the local max producer offset from disk
	ProducerOffsetLocal(topic *TopicDataId) (int64, error)

	// Get the lowest offset value of any group for a given topic+token+index
	MinOffset(topic string, token Token, index RangeIndex) *Offset
}
