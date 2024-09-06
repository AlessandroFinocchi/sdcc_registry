package services

import (
	"context"
	"fmt"
	"github.com/AlessandroFinocchi/sdcc_common/pb"
	"github.com/AlessandroFinocchi/sdcc_common/utils"
	"log"
	"math/rand"
	"os"
	m "sdcc_registry/model"
	ur "sdcc_registry/utils"
	"strconv"
	"sync"
	"time"
)

type Connector struct {
	pb.UnimplementedConnectorServer
	Mutex  *sync.Mutex
	NodesW *m.NodeListWrapper
	c      uint64
	logger ur.MyLogger
}

func NewConnector(mutex *sync.Mutex, nodesW *m.NodeListWrapper) *Connector {
	membershipNodesNumber, err := utils.ReadConfigUInt64("config.ini", "membership", "c")
	logging, errL := strconv.ParseBool(os.Getenv(ur.LoggingEnv))
	if err != nil || errL != nil {
		log.Fatalf("Could not read configuration in connector: %v", err)
	}

	return &Connector{
		Mutex:  mutex,
		NodesW: nodesW,
		c:      membershipNodesNumber,
		logger: ur.NewMyLogger(logging),
	}
}

func (c *Connector) getNodes(ctx context.Context) (*pb.NodeList, error) {

	if err := utils.ContextError(ctx); err != nil {
		return nil, err
	}

	sendingNodes := make([]*pb.Node, len(c.NodesW.NodeList))
	copy(sendingNodes, c.NodesW.NodeList)

	rand.NewSource(time.Now().Unix())
	rand.Shuffle(len(sendingNodes), func(i, j int) {
		sendingNodes[i], sendingNodes[j] = sendingNodes[j], sendingNodes[i]
	})

	lastIndex := min(c.c, uint64(len(sendingNodes)))

	return &pb.NodeList{Nodes: sendingNodes[:lastIndex]}, nil
}

func (c *Connector) Connect(ctx context.Context, in *pb.Node) (*pb.NodeList, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if err := utils.ContextError(ctx); err != nil {
		return nil, err
	}

	sendingNodes, err := c.getNodes(ctx)

	if err != nil {
		c.logger.Log(fmt.Sprintf("Could not get nodes to send: %v\n", err))
	}

	c.NodesW.Add(in)

	c.logger.Log(fmt.Sprintf("Connected to %s:%s:%d", in.GetId(), in.GetMembershipIp(), in.GetMembershipPort()))
	return sendingNodes, nil
}

func (c *Connector) Disconnect(ctx context.Context, in *pb.Node) (*pb.Empty, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if err := utils.ContextError(ctx); err != nil {
		return nil, err
	}

	c.NodesW.Remove(in.GetId())

	c.logger.Log(fmt.Sprintf("Disconnected from %s:%s:%d", in.GetId(), in.GetMembershipIp(), in.GetMembershipPort()))
	return nil, nil
}
