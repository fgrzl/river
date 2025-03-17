package server

import (
	"encoding/json"
	"fmt"

	"github.com/fgrzl/lexkey"
	"github.com/google/uuid"
)

type GetStatus struct{}

func (g *GetStatus) GetRoute() string {
	return "get_node_count"
}

type ClusterStatus struct {
	NodeCount int `json:"node_count"`
}

type Peek struct {
	Space   string `json:"space"`
	Segment string `json:"segment"`
}

func (g *Peek) GetRoute() string {
	return "peek"
}

type Produce struct {
	Space   string `json:"space"`
	Segment string `json:"segment"`
}

func (g *Produce) GetRoute() string {
	return "produce"
}

type ConsumeSpace struct {
	Space        string `json:"space"`
	MinTimestamp int64  `json:"min_timestamp"`
	MaxTimestamp int64  `json:"max_timestamp"`
	Offset       []byte `json:"offset"`
}

func (g *ConsumeSpace) GetRoute() string {
	return "consume_space"
}

type ConsumeSegment struct {
	Space   string `json:"space"`
	Segment string `json:"segment"`

	// The minimum sequence number to consume.
	MinSequence  uint64 `json:"min_sequence"`
	MinTimestamp int64  `json:"min_timestamp"`
	MaxSequence  uint64 `json:"max_sequence"`
	MaxTimestamp int64  `json:"max_timestamp"`
}

func (g *ConsumeSegment) GetRoute() string {
	return "consume_segment"
}

//
// Data Management
//

type GetSpaces struct{}

func (g *GetSpaces) GetRoute() string {
	return "get_spaces"
}

type GetSegments struct {
	Space string `json:"space"`
}

func (g *GetSegments) GetRoute() string {
	return "get_segments"
}

type EnumerateSpace struct {
	Space        string `json:"space"`
	MinTimestamp int64  `json:"min_timestamp"`
	MaxTimestamp int64  `json:"max_timestamp"`
	Offset       []byte `json:"offset"`
}

func (g *EnumerateSpace) GetRoute() string {
	return "enumerate_space"
}

type EnumerateSegment struct {
	Space   string `json:"space"`
	Segment string `json:"segment"`

	// The minimum sequence number to consume.
	MinSequence  uint64 `json:"min_sequence"`
	MinTimestamp int64  `json:"min_timestamp"`
	MaxSequence  uint64 `json:"max_sequence"`
	MaxTimestamp int64  `json:"max_timestamp"`
}

func (g *EnumerateSegment) GetRoute() string {
	return "enumerate_segment"
}

type CheckSpaceOffset struct {
	ID     uuid.UUID     `json:"id"`
	Node   uuid.UUID     `json:"node"`
	Space  string        `json:"space"`
	Offset lexkey.LexKey `json:"offset"`
}

func (c *CheckSpaceOffset) GetRoute() string {
	return "check_space_offset"
}

func (c *CheckSpaceOffset) ToACK(node uuid.UUID) *ACK {
	return &ACK{
		ID:   c.ID,
		Node: node,
	}
}

func (c *CheckSpaceOffset) ToNACK(node uuid.UUID) *NACK {
	return &NACK{
		ID:   c.ID,
		Node: node,
	}
}

type CheckSegmentOffset struct {
	ID      uuid.UUID     `json:"id"`
	Node    uuid.UUID     `json:"node"`
	Space   string        `json:"space"`
	Segment string        `json:"segment"`
	Offset  lexkey.LexKey `json:"offset"`
}

func (c *CheckSegmentOffset) GetRoute() string {
	return "check_segment_offset"
}

func (c *CheckSegmentOffset) ToACK(node uuid.UUID) *ACK {
	return &ACK{
		ID:   c.ID,
		Node: node,
	}
}

func (c *CheckSegmentOffset) ToNACK(node uuid.UUID) *NACK {
	return &NACK{
		ID:   c.ID,
		Node: node,
	}
}

//
// Transaction Management
//

const (
	UNCOMMITTED = "uncommitted"
	COMMITTED   = "committed"
	FINALIZED   = "finalized"
)

type Transaction struct {
	TRX           TRX      `json:"trx"`
	Space         string   `json:"space"`
	Segment       string   `json:"segment"`
	FirstSequence uint64   `json:"first_sequence"`
	LastSequence  uint64   `json:"last_sequence"`
	Entries       []*Entry `json:"entries"`
	Timestamp     int64    `json:"timestamp"`
}

func (a *Transaction) GetRoute() string {
	return fmt.Sprintf("%T", a)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {

	raw, err := EncodeTransaction(t)
	if err != nil {
		return nil, err
	}

	wrapper := struct {
		D []byte `json:"d"`
	}{D: raw}

	return json.Marshal(wrapper)
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	wrapper := struct {
		D []byte `json:"d"`
	}{}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return fmt.Errorf("failed to unmarshal wrapper: %w", err)
	}

	if len(wrapper.D) == 0 {
		return fmt.Errorf("compressed data is empty")
	}

	return DecodeTransaction(wrapper.D, t)
}

type TRX struct {
	ID     uuid.UUID `json:"id"`
	Node   uuid.UUID `json:"node"`
	Number uint64    `json:"number"`
}

func (a *TRX) GetRoute() string {
	return fmt.Sprintf("%T.%v", a, a.ID)
}

func (a *TRX) ToACK(node uuid.UUID) *ACK {
	return &ACK{
		ID:   a.ID,
		Node: node,
	}
}

func (a *TRX) ToNACK(node uuid.UUID) *NACK {
	return &NACK{
		ID:   a.ID,
		Node: node,
	}
}

type Commit struct {
	TRX     TRX    `json:"trx"`
	Space   string `json:"space"`
	Segment string `json:"segment"`
}

func (a *Commit) GetRoute() string {
	return "trx.commit"
}

type Reconcile struct {
	TRX     TRX    `json:"trx"`
	Space   string `json:"space"`
	Segment string `json:"segment"`
}

func (a *Reconcile) GetRoute() string {
	return "trx.reconcile"
}

type Rollback struct {
	TRX     TRX    `json:"trx"`
	Space   string `json:"space"`
	Segment string `json:"segment"`
}

func (a *Rollback) GetRoute() string {
	return "trx.rollback"
}

//
// Node Management
//

type Synchronize struct {
	OffsetsBySpace map[string]lexkey.LexKey `json:"offsets_by_space"`
}

func (a *Synchronize) GetRoute() string {
	return "node.synchronize"
}

// QuorumChanged represents a quorum update
type QuorumChanged struct {
	Node uuid.UUID `json:"node"`
}

func (q QuorumChanged) GetRoute() string {
	return "node.quorum_changed"
}

// NodeHeartbeat represents a node failure event
type NodeHeartbeat struct {
	Node uuid.UUID `json:"node"`
}

func (h *NodeHeartbeat) GetRoute() string {
	return "node.healthcheck"
}

// NodeShutdown notifies that a node has gone down
type NodeShutdown struct {
	Node uuid.UUID `json:"node"`
}

func (n *NodeShutdown) GetRoute() string {
	return "node.shutdown"
}

//
// ACK and NACK
//

type ACK struct {
	ID   uuid.UUID `json:"id"`
	Node uuid.UUID `json:"node"`
}

func (a *ACK) GetRoute() string {
	return GetReplyRoute(a.ID)
}

type NACK struct {
	ID   uuid.UUID `json:"id"`
	Node uuid.UUID `json:"node"`
}

func (a *NACK) GetRoute() string {
	return GetReplyRoute(a.ID)
}

func GetReplyRoute(messageId uuid.UUID) string {
	return "reply." + messageId.String()
}
