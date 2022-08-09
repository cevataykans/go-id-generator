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

	custom_epoch int64
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

		custom_epoch: epoch.UnixMilli(),
		nodeId:       node,

		counter:   0,
		timestamp: time.Now().UnixMilli(),
	}, nil
}

func (n *FlakeNode) GenerateToChannel(c chan<- SnowFlakeID) {

	msTicker := time.NewTicker(time.Millisecond)
	now := time.Now().UnixMilli()
	for {

		select {
		case <-msTicker.C:

			now++
			n.counter = 0

		case c <- SnowFlakeID(((now - n.custom_epoch) << TimeShift) |
			(n.nodeId << NodeShift) |
			n.counter):

			n.counter = (n.counter + 1) & MaxCounter
			if n.counter == 0 {
				time.Sleep(time.Millisecond)
			}
		}
	}
}

func (n *FlakeNode) GenerateToChannelByReq(c <-chan *FlakeRequest) {

	msTicker := time.NewTicker(time.Millisecond)
	now := time.Now().UnixMilli()
	for {

		select {
		case <-msTicker.C:

			now++
			n.counter = 0

		case req := <-c:

			req.id = SnowFlakeID(
				((now - n.custom_epoch) << TimeShift) |
					(n.nodeId << NodeShift) |
					n.counter)
			close(req.done)

			n.counter = (n.counter + 1) & MaxCounter
			if n.counter == 0 {
				time.Sleep(time.Millisecond)
			}
		}
	}
}

func (n *FlakeNode) GeneratePerRequest(c <-chan *FlakeRequest) {

	for req := range c {

		now := time.Now().UnixMilli()
		id := SnowFlakeID((now-n.custom_epoch)<<TimeShift |
			(n.nodeId << NodeShift) | n.counter)

		if now == n.timestamp {
			n.counter = (n.counter + 1) / MaxCounter
			if n.counter == 0 {
				//TODO
				fmt.Println("244 ns goal reached!")
			}
		} else {
			n.timestamp = now
			n.counter = 0
		}

		req.id = id
		close(req.done)
	}
}

func (n *FlakeNode) GenerateByLock() SnowFlakeID {

	n.Lock()

	now := time.Now().UnixMilli()
	id := SnowFlakeID((now-n.custom_epoch)<<TimeShift |
		(n.nodeId << NodeShift) | n.counter)

	if now == n.timestamp {
		n.counter = (n.counter + 1) & MaxCounter
		if n.counter == 0 {
			//TODO
			fmt.Println("244 ns goal reached!")
		}
	} else {
		n.timestamp = now
		n.counter = 0
	}
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
