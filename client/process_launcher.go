package client

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/longhorn/longhorn-engine-launcher/api"
	"github.com/longhorn/longhorn-engine-launcher/rpc"
	"github.com/longhorn/longhorn-engine-launcher/types"
)

type ProcessLauncherClient struct {
	Address string
}

func NewProcessLauncherClient(address string) *ProcessLauncherClient {
	return &ProcessLauncherClient{
		Address: address,
	}
}

func RPCToProcess(obj *rpc.ProcessResponse) *api.Process {
	return &api.Process{
		Name:          obj.Spec.Name,
		Binary:        obj.Spec.Binary,
		Args:          obj.Spec.Args,
		PortCount:     obj.Spec.PortCount,
		PortArgs:      obj.Spec.PortArgs,
		ProcessStatus: RPCToProcessStatus(obj.Status),
	}
}

func RPCToProcessList(obj *rpc.ProcessListResponse) map[string]*api.Process {
	ret := map[string]*api.Process{}
	for name, p := range obj.Processes {
		ret[name] = RPCToProcess(p)
	}
	return ret
}

func (cli *ProcessLauncherClient) ProcessCreate(name, binary string, portCount int, args, portArgs []string) (*api.Process, error) {
	if name == "" || binary == "" {
		return nil, fmt.Errorf("failed to start process: missing required parameter")
	}

	conn, err := grpc.Dial(cli.Address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("cannot connect process launcher service to %v: %v", cli.Address, err)
	}
	defer conn.Close()

	client := rpc.NewLonghornProcessLauncherServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), types.GRPCServiceTimeout)
	defer cancel()

	p, err := client.ProcessCreate(ctx, &rpc.ProcessCreateRequest{
		Spec: &rpc.ProcessSpec{
			Name:      name,
			Binary:    binary,
			Args:      args,
			PortCount: int32(portCount),
			PortArgs:  portArgs,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start process: %v", err)
	}
	return RPCToProcess(p), nil
}

func (cli *ProcessLauncherClient) ProcessDelete(name string) (*api.Process, error) {
	if name == "" {
		return nil, fmt.Errorf("failed to delete process: missing required parameter name")
	}

	conn, err := grpc.Dial(cli.Address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("cannot connect process launcher service to %v: %v", cli.Address, err)
	}
	defer conn.Close()

	client := rpc.NewLonghornProcessLauncherServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), types.GRPCServiceTimeout)
	defer cancel()

	p, err := client.ProcessDelete(ctx, &rpc.ProcessDeleteRequest{
		Name: name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete process %v: %v", name, err)
	}
	return RPCToProcess(p), nil
}

func (cli *ProcessLauncherClient) ProcessGet(name string) (*api.Process, error) {
	if name == "" {
		return nil, fmt.Errorf("failed to get process: missing required parameter name")
	}

	conn, err := grpc.Dial(cli.Address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("cannot connect process launcher service to %v: %v", cli.Address, err)
	}
	defer conn.Close()

	client := rpc.NewLonghornProcessLauncherServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), types.GRPCServiceTimeout)
	defer cancel()

	p, err := client.ProcessGet(ctx, &rpc.ProcessGetRequest{
		Name: name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get process %v: %v", name, err)
	}
	return RPCToProcess(p), nil
}

func (cli *ProcessLauncherClient) ProcessList() (map[string]*api.Process, error) {
	conn, err := grpc.Dial(cli.Address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("cannot connect process launcher service to %v: %v", cli.Address, err)
	}
	defer conn.Close()

	client := rpc.NewLonghornProcessLauncherServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), types.GRPCServiceTimeout)
	defer cancel()

	ps, err := client.ProcessList(ctx, &rpc.ProcessListRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list processes: %v", err)
	}
	return RPCToProcessList(ps), nil
}