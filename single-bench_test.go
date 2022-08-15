package main

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkLock(b *testing.B) {

	idGen, _ := NewFlakeNode(time.Now().Round(time.Millisecond), 0)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		idGen.GenerateByLock()
	}
}

func BenchmarkToChannel(b *testing.B) {

	idGen, _ := NewFlakeNode(time.Now().Round(time.Millisecond), 0)
	flakeChan := make(chan SnowFlakeID)
	go idGen.GenerateToChannel(flakeChan)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		<-flakeChan
	}
}

func BenchmarkChannelReq(b *testing.B) {

	for _, inputSize := range []int{0, 1, 4, 16, 256, 1024} {
		cap := inputSize
		b.Run(
			fmt.Sprintf("inputSize= %v", inputSize), func(b *testing.B) {

				idGen, _ := NewFlakeNode(time.Now().Round(time.Millisecond), 0)
				b.ReportAllocs()

				flakeChan := make(chan *FlakeRequest, cap)
				go idGen.GenerateToChannelByReq(flakeChan)

				for i := 0; i < b.N; i++ {
					req := &FlakeRequest{done: make(chan struct{})}
					flakeChan <- req
					<-req.done
				}
			},
		)
	}
}

func BenchmarkChannelReqAutoTime(b *testing.B) {

	for _, inputSize := range []int{0, 1, 4, 16, 256, 1024} {
		cap := inputSize
		b.Run(
			fmt.Sprintf("inputSize= %v", inputSize), func(b *testing.B) {

				idGen, _ := NewFlakeNode(time.Now().Round(time.Millisecond), 0)
				b.ReportAllocs()

				flakeChan := make(chan *FlakeRequest, cap)
				go idGen.GenerateToChannelByReqAutoTime(flakeChan)

				for i := 0; i < b.N; i++ {
					req := &FlakeRequest{done: make(chan struct{})}
					flakeChan <- req
					<-req.done
				}
			},
		)
	}
}
