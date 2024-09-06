package services

import (
	"context"
	"fmt"
	"github.com/AlessandroFinocchi/sdcc_common/pb"
	"github.com/AlessandroFinocchi/sdcc_common/utils"
	"log"
	"os"
	m "sdcc_registry/model"
	ur "sdcc_registry/utils"
	"strconv"
	"sync"
)

type Heartbeat struct {
	pb.UnimplementedHeartbeatServer
	Mutex           *sync.Mutex
	NodesW          *m.NodeListWrapper
	cycles          uint64
	cyclesThreshold uint64
	logger          ur.MyLogger
}

func NewHeartbeat(mutex *sync.Mutex, nodesW *m.NodeListWrapper) *Heartbeat {
	cyclesThreshold, err := utils.ReadConfigUInt64("config.ini", "heartbeat", "cycles_threshold")
	logging, errL := strconv.ParseBool(os.Getenv(ur.LoggingEnv))
	if err != nil || errL != nil {
		log.Fatalf("Could not read configuration in heartbeat: %v", err)
	}
	return &Heartbeat{
		Mutex:           mutex,
		NodesW:          nodesW,
		cycles:          0,
		cyclesThreshold: cyclesThreshold,
		logger:          ur.NewMyLogger(logging),
	}
}

func (h *Heartbeat) Beat(ctx context.Context, node *pb.Node) (*pb.Empty, error) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	if err := utils.ContextError(ctx); err != nil {
		return nil, err
	}

	h.logger.Log(fmt.Sprintf("Beat from node %s", node.GetId()))

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
			h.logger.Log(fmt.Sprintf("Node %s has timed out.", node.GetId()))
			h.NodesW.Remove(node.GetId())
		}
	}

	h.logger.Log("Remaining nodes:")
	for _, node := range h.NodesW.NodeList {
		h.logger.Log(fmt.Sprintf("%s: %s:%d", node.GetId(), node.GetMembershipIp(), node.GetMembershipPort()))
	}
	h.logger.Log("")
}
