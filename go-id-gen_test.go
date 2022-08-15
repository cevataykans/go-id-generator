package main

import (
	"testing"
	"time"
)

func TestGeneratorToChannel(t *testing.T) {

	epoch := time.Now().Round(time.Millisecond)
	gen, err := NewFlakeNode(epoch, 0)
	if err != nil {
		t.Fatal(err)
	}

	flakeChan := make(chan SnowFlakeID)
	go gen.GenerateToChannel(flakeChan)

	reqTime := time.Now()
	id := <-flakeChan
	t.Logf("Id: %v", id)

	if id.GetCounter() != 0 {
		t.Fatalf("Expected counter: %v; actual counter: %v", 0, id.GetCounter())
	}

	if id.GetNode() != 0 {
		t.Fatalf("expected node id: %v; actual: %v", 0, id.GetNode())
	}

	idTime := id.GetTime(epoch).Round(time.Millisecond)
	if !idTime.Equal(epoch) && idTime.Add(time.Millisecond).Equal(reqTime) {
		t.Fatalf("expected time: %v; actual creation time: %v\n", reqTime, idTime)
	}

	idSet := make(map[SnowFlakeID]bool)
	for i := 0; i < 1_000_000; i++ {
		id = <-flakeChan
		if idSet[id] == true {
			t.Fatal("duplicate id generated")
		}
		idSet[id] = true
	}
}

func TestGeneratorWithLock(t *testing.T) {

	epoch := time.Now().Round(time.Millisecond)
	gen, err := NewFlakeNode(epoch, 0)
	if err != nil {
		t.Fatal(err)
	}

	reqTime := time.Now()
	id := gen.GenerateByLock()
	t.Logf("Id: %v", id)

	if id.GetCounter() != 0 {
		t.Fatalf("Expected counter: %v; actual counter: %v", 0, id.GetCounter())
	}

	if id.GetNode() != 0 {
		t.Fatalf("expected node id: %v; actual: %v", 0, id.GetNode())
	}

	idTime := id.GetTime(epoch).Round(time.Millisecond)
	if !idTime.Equal(epoch) && idTime.Add(time.Millisecond).Equal(reqTime) {
		t.Fatalf("expected time: %v; actual creation time: %v\n", reqTime, idTime)
	}

	idSet := make(map[SnowFlakeID]bool)
	for i := 0; i < 1_000_000; i++ {
		id = gen.GenerateByLock()
		if idSet[id] == true {
			t.Fatal("duplicate id generated")
		}
		idSet[id] = true
	}
}

func TestGeneratorFromReqChannel(t *testing.T) {

	epoch := time.Now().Round(time.Millisecond)
	gen, err := NewFlakeNode(epoch, 0)
	if err != nil {
		t.Fatal(err)
	}

	reqTime := time.Now()
	flakeChan := make(chan *FlakeRequest)
	go gen.GenerateToChannelByReq(flakeChan)

	flakeReq := FlakeRequest{done: make(chan struct{})}
	flakeChan <- &flakeReq
	<-flakeReq.done
	id := flakeReq.id

	t.Logf("Id: %v", id)

	if id.GetCounter() != 0 {
		t.Fatalf("Expected counter: %v; actual counter: %v", 0, id.GetCounter())
	}

	if id.GetNode() != 0 {
		t.Fatalf("expected node id: %v; actual: %v", 0, id.GetNode())
	}

	idTime := id.GetTime(epoch).Round(time.Millisecond)
	if !idTime.Equal(epoch) && idTime.Add(time.Millisecond).Equal(reqTime) {
		t.Fatalf("expected time: %v; actual creation time: %v\n", reqTime, idTime)
	}

	idSet := make(map[SnowFlakeID]bool)
	for i := 0; i < 1_000_000; i++ {

		flakeReq = FlakeRequest{done: make(chan struct{})}
		flakeChan <- &flakeReq
		<-flakeReq.done
		id = flakeReq.id

		if idSet[id] == true {
			t.Fatal("duplicate id generated")
		}
		idSet[id] = true
	}
}

func TestGeneratorFromReqChannelWithAutoTime(t *testing.T) {

	epoch := time.Now().Round(time.Millisecond)
	gen, err := NewFlakeNode(epoch, 0)
	if err != nil {
		t.Fatal(err)
	}

	reqTime := time.Now()
	flakeChan := make(chan *FlakeRequest)
	go gen.GenerateToChannelByReqAutoTime(flakeChan)

	flakeReq := FlakeRequest{done: make(chan struct{})}
	flakeChan <- &flakeReq
	<-flakeReq.done
	id := flakeReq.id

	t.Logf("Id: %v", id)

	if id.GetCounter() != 0 {
		t.Fatalf("Expected counter: %v; actual counter: %v", 0, id.GetCounter())
	}

	if id.GetNode() != 0 {
		t.Fatalf("expected node id: %v; actual: %v", 0, id.GetNode())
	}

	idTime := id.GetTime(epoch).Round(time.Millisecond)
	if !idTime.Equal(epoch) && idTime.Add(time.Millisecond).Equal(reqTime) {
		t.Fatalf("expected time: %v; actual creation time: %v\n", reqTime, idTime)
	}

	idSet := make(map[SnowFlakeID]bool)
	for i := 0; i < 1_000_000; i++ {

		flakeReq = FlakeRequest{done: make(chan struct{})}
		flakeChan <- &flakeReq
		<-flakeReq.done
		id = flakeReq.id

		if idSet[id] == true {
			t.Fatal("duplicate id generated")
		}
		idSet[id] = true
	}
}
