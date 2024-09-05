package services

import (
	"context"
	"fmt"
	"github.com/AlessandroFinocchi/sdcc_common/pb"
	"github.com/AlessandroFinocchi/sdcc_common/utils"
	"log"
	m "sdcc_registry/model"
	"sync"
)

type Heartbeat struct {
	pb.UnimplementedHeartbeatServer
	Mutex           *sync.Mutex
	NodesW          *m.NodeListWrapper
	cycles          uint64
	cyclesThreshold uint64
}

func NewHeartbeat(mutex *sync.Mutex, nodesW *m.NodeListWrapper) *Heartbeat {
	cyclesThreshold, err := utils.ReadConfigUInt64("config.ini", "heartbeat", "cycles_threshold")
	if err != nil {
		log.Fatalf("Could not read cycles_threshold from config file: %v", err)
	}
	return &Heartbeat{
		Mutex:           mutex,
		NodesW:          nodesW,
		cycles:          0,
		cyclesThreshold: cyclesThreshold,
	}
}

func (h *Heartbeat) Beat(ctx context.Context, node *pb.Node) (*pb.Empty, error) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	if err := utils.ContextError(ctx); err != nil {
		return nil, err
	}

	fmt.Println("Beat from node " + node.GetId())

	if _, ok := h.NodesW.NodeMap[node.GetId()]; !ok {
		h.NodesW.Add(node)
	}

	h.NodesW.NodeMap[node.GetId()].Tick++
	h.NodesW.NodeMap[node.GetId()].LastUpdated = h.cycles

	return nil, nil
}

func (h *Heartbeat) OnTimeout() {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	h.cycles++

	for _, node := range h.NodesW.NodeMap {
		if h.cycles-node.LastUpdated > h.cyclesThreshold {
			fmt.Println("Node " + node.GetId() + " has timed out.")
			h.NodesW.Remove(node.GetId())
		}
	}

	fmt.Println("Remaining nodes:")
	for _, node := range h.NodesW.NodeList {
		fmt.Println(node.GetId(), ": ", node.GetMembershipIp(), ":", node.GetMembershipPort())
	}
	fmt.Println()
}
