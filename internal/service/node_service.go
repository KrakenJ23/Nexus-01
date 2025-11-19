package service

import (
	"Nexus/internal/adapters/runtime"
	"Nexus/internal/core"
	"Nexus/internal/ports"
	"fmt"
)

// NodeService is the kind of mastermind of the orchestration of the nodes
type NodeService struct {
	runtime ports.ContainerRuntime
}

// NewNodeService is the method that will allow us to create from the logic metier a new container
func NewNodeService() (*NodeService, error) {
	// Here, we are going to instantiate the runtime adapter
	rt, runtimeInitializationError := runtime.NewLibContainerRuntime()
	if runtimeInitializationError != nil {
		return nil, fmt.Errorf("unable to initialize libcontainer runtime %v", runtimeInitializationError.Error())
	}
	return &NodeService{runtime: rt}, nil
}

// Now let's deal with the application logic to create a new node
func (s *NodeService) CreateNode(name string, mem int64, cpuShare uint64) (*core.NodeState, error) {
	// Let's perform a little validation
	if name == "" {
		return nil, fmt.Errorf("the node must have a name")
	}

	// let's now create the configuration to be launched
	rootfsPath := "/var/lib/nexus/images/alpine-base" // this the location of the root filesystem
	conf := core.NodeConfig{
		ID:         name,
		Hostname:   name,
		RootfsPath: rootfsPath,
		Memory:     mem,
		CPUShares:  cpuShare,
		Command:    []string{"/bin/bash", "-c", "sleep 3600"}, // our process
	}

	// Now let's call the adapter by using the interface
	state, err := s.runtime.CreateAndStart(conf)
	if err != nil {
		return nil, fmt.Errorf("unable to create a new node %v", err.Error())
	}
	return state, nil
}
