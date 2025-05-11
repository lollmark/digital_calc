package tests

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/lollmark/calculator_go/internal"
	"github.com/lollmark/calculator_go/proto/calc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeServer struct {
	calc.UnimplementedCalcServer
	taskCalled       bool
	postResultCalled bool
	task             *calc.TaskResp
	resultReq        *calc.ResultReq
}

func (f *fakeServer) GetTask(ctx context.Context, _ *calc.Empty) (*calc.TaskResp, error) {
	f.taskCalled = true
	if f.task == nil {
		return nil, status.Error(codes.NotFound, "no task")
	}
	return f.task, nil
}

func (f *fakeServer) PostResult(ctx context.Context, in *calc.ResultReq) (*calc.Empty, error) {
	f.postResultCalled = true
	f.resultReq = in
	return &calc.Empty{}, nil
}

func TestAgent_WorkerFlow(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	fake := &fakeServer{
		task: &calc.TaskResp{
			Id:            "task1",
			Arg1:          2,
			Arg2:          3,
			Operation:     "*",
			OperationTime: 10, 
		},
	}
	calc.RegisterCalcServer(grpcServer, fake)
	go grpcServer.Serve(lis)
	defer grpcServer.Stop()

	existing := setEnv("ORCHESTRATOR_URL", "http://"+lis.Addr().String())
	defer restoreEnv("ORCHESTRATOR_URL", existing)

	agent := application.NewAgent()
	agent.ComputingPower = 1

	done := make(chan struct{})
	go func() {
		agent.Worker(0)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)

	if !fake.taskCalled {
		t.Error("expected GetTask to be called")
	}
	if !fake.postResultCalled {
		t.Error("expected PostResult to be called")
	}

	if fake.resultReq.Id != "task1" {
		t.Errorf("expected result id 'task1', got %s", fake.resultReq.Id)
	}
	if fake.resultReq.Result != 6 {
		t.Errorf("expected result 6, got %f", fake.resultReq.Result)
	}
}

func setEnv(key, val string) string {
	old := os.Getenv(key)
	os.Setenv(key, val)
	return old
}

func restoreEnv(key, val string) {
	if val == "" {
		os.Unsetenv(key)
	} else {
		os.Setenv(key, val)
	}
}
