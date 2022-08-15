package main

import (
	"fmt"
	"sync"
	"time"
)

type SnowFlakeID uint64

const (
	NodeBits    = 10
	CounterBits = 12
	TimeBits    = 41

	TimeShift = NodeBits + CounterBits
	NodeShift = CounterBits

	MaxCounter   = (1 << CounterBits) - 1
	MaxNodeCount = (1 << NodeBits) - 1
)

type FlakeNode struct {
	sync.Mutex

	custom_epoch time.Time
	nodeId       int64

	counter   int64
	timestamp int64
}

type FlakeRequest struct {
	id   SnowFlakeID
	done chan struct{}
}

func NewFlakeNode(epoch time.Time, node int64) (*FlakeNode, error) {

	if node > MaxNodeCount {
		return nil, fmt.Errorf("max node count reached")
	}

	return &FlakeNode{

		custom_epoch: epoch,
		nodeId:       node,

		counter:   0,
		timestamp: epoch.UnixMilli(),
	}, nil
}

func (n *FlakeNode) GenerateToChannel(c chan<- SnowFlakeID) {

	msTicker := time.NewTicker(time.Millisecond)
	now := time.Now().UnixMilli()
	epoch := n.custom_epoch.UnixMilli()
	for {

		select {
		case <-msTicker.C:

			now++
			n.counter = 0

		case c <- SnowFlakeID(((now - epoch) << TimeShift) |
			(n.nodeId << NodeShift) |
			n.counter):

			n.counter = (n.counter + 1) & MaxCounter
			if n.counter == 0 {
				<-msTicker.C
				now++
				n.counter = 0
			}
		}
	}
}

func (n *FlakeNode) GenerateToChannelByReqAutoTime(c <-chan *FlakeRequest) {

	msTicker := time.NewTicker(time.Millisecond)
	now := time.Now().UnixMilli()
	epoch := n.custom_epoch.UnixMilli()
	for {

		select {
		case <-msTicker.C:

			now++
			n.counter = 0

		case req := <-c:

			req.id = SnowFlakeID(
				((now - epoch) << TimeShift) |
					(n.nodeId << NodeShift) |
					n.counter)

			n.counter = (n.counter + 1) & MaxCounter
			if n.counter == 0 {
				<-msTicker.C
				now++
				n.counter = 0
			}
			close(req.done)
		}
	}
}

func (n *FlakeNode) GenerateToChannelByReq(c <-chan *FlakeRequest) {

	for req := range c {

		now := time.Since(n.custom_epoch).Milliseconds()

		if now == n.timestamp {
			n.counter = (n.counter + 1) & MaxCounter
			if n.counter == 0 {
				for now <= n.timestamp {
					now = time.Since(n.custom_epoch).Milliseconds()
				}
			}
		} else {
			n.counter = 0
		}
		n.timestamp = now

		req.id = SnowFlakeID(now<<TimeShift |
			(n.nodeId << NodeShift) | n.counter)

		close(req.done)
	}
}

func (n *FlakeNode) GenerateByLock() SnowFlakeID {

	n.Lock()

	now := time.Since(n.custom_epoch).Milliseconds()
	if now == n.timestamp {
		n.counter = (n.counter + 1) & MaxCounter
		if n.counter == 0 {
			for now <= n.timestamp {
				now = time.Since(n.custom_epoch).Milliseconds()
			}
		}
	} else {
		n.counter = 0
	}
	n.timestamp = now

	id := SnowFlakeID((now << TimeShift) |
		(n.nodeId << NodeShift) | n.counter)

	n.Unlock()
	return id
}

func (id SnowFlakeID) GetNode() int64 {

	return (int64(id) & (MaxNodeCount << NodeShift)) >> NodeShift
}

func (id SnowFlakeID) GetCounter() int64 {

	return int64(id) & MaxCounter
}

// Time since custom epoch
func (id SnowFlakeID) GetTime(epoch time.Time) time.Time {

	timestamp := (int64(id) >> TimeShift) + epoch.UnixMilli()
	return time.UnixMilli(timestamp)
}
