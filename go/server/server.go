package main

import (
	"context"
	"time"

	"github.com/ktny/study-grpc/go/deepthought"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ComputeServerを実装する型
type Server struct {
	deepthought.UnimplementedComputeServer
}

// インタフェースが実装できていることをコンパイル時に確認するおまじない
var _ deepthought.ComputeServer = &Server{}

// Boot RPCの型
func (s *Server) Boot(req *deepthought.BootRequest, stream deepthought.Compute_BootServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-time.After(1 * time.Second):
		}

		if err := stream.Send(&deepthought.BootResponse{
			Message: "I think therefore I am.",
		}); err != nil {
			return err
		}
	}
}

// Infer RPCの型
func (s *Server) Infer(ctx context.Context, req *deepthought.InferRequest) (*deepthought.InferResponse, error) {
	switch req.Query {
	case "Life", "Universe", "Everything":
	default:
		return nil, status.Error(codes.InvalidArgument, "Contemplate your query")
	}

	deadline, ok := ctx.Deadline()

	// 指定されていない、もしくは十分な時間があれば回答
	if !ok || time.Until(deadline) > 750*time.Millisecond {
		time.Sleep(750 * time.Millisecond)
		return &deepthought.InferResponse{
			Answer:      42,
			Description: []string{"I checked it"},
		}, nil
	}

	return nil, status.Error(codes.DeadlineExceeded, "It would take longer")
}
