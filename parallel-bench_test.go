package main

import (
	"fmt"
	"testing"
	"time"
)

func Benchmark_LockParallel(b *testing.B) {

	idGen, err := NewFlakeNode(time.Now().Round(time.Millisecond), 0)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			idGen.GenerateByLock()
		}
	})
}

func Benchmark_ToChannelParallel(b *testing.B) {

	idGen, err := NewFlakeNode(time.Now().Round(time.Millisecond), 0)
	if err != nil {
		b.Fatal(err)
	}

	flakeChan := make(chan SnowFlakeID)
	go idGen.GenerateToChannel(flakeChan)

	b.ReportAllocs()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			<-flakeChan
		}
	})
}

func Benchmark_ChannelReqParallel(b *testing.B) {

	for _, inputSize := range []int{0, 1, 4, 16, 256, 1024} {
		cap := inputSize
		b.Run(
			fmt.Sprintf("inputSize= %v", inputSize), func(b *testing.B) {

				idGen, err := NewFlakeNode(time.Now().Round(time.Millisecond), 0)
				if err != nil {
					b.Fatal(err)
				}

				flakeChan := make(chan *FlakeRequest, cap)
				go idGen.GenerateToChannelByReq(flakeChan)

				b.ReportAllocs()
				b.RunParallel(func(p *testing.PB) {
					for p.Next() {
						req := &FlakeRequest{done: make(chan struct{})}
						flakeChan <- req
						<-req.done
					}
				})
			},
		)
	}
}

func Benchmark_ChannelReqAutoTimeParallel(b *testing.B) {

	for _, inputSize := range []int{0, 1, 4, 16, 256, 1024} {
		cap := inputSize
		b.Run(
			fmt.Sprintf("inputSize= %v", inputSize), func(b *testing.B) {

				idGen, err := NewFlakeNode(time.Now().Round(time.Millisecond), 0)
				if err != nil {
					b.Fatal(err)
				}

				flakeChan := make(chan *FlakeRequest, cap)
				go idGen.GenerateToChannelByReqAutoTime(flakeChan)

				b.ReportAllocs()
				b.RunParallel(func(p *testing.PB) {
					for p.Next() {
						req := &FlakeRequest{done: make(chan struct{})}
						flakeChan <- req
						<-req.done
					}
				})
			},
		)
	}
}
