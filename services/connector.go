package services

import (
	"context"
	"fmt"
	"github.com/AlessandroFinocchi/sdcc_common/pb"
	"github.com/AlessandroFinocchi/sdcc_common/utils"
	"log"
	"math/rand"
	m "sdcc_registry/model"
	"sync"
	"time"
)

type Connector struct {
	pb.UnimplementedConnectorServer
	Mutex  *sync.Mutex
	NodesW *m.NodeListWrapper
	c      uint64
}

func NewConnector(mutex *sync.Mutex, nodesW *m.NodeListWrapper) *Connector {
	membershipNodesNumber, err := utils.ReadConfigUInt64("config.ini", "membership", "c")
	if err != nil {
		log.Fatalf("Could not read c from config file: %v", err)
	}

	return &Connector{
		Mutex:  mutex,
		NodesW: nodesW,
		c:      membershipNodesNumber,
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
		fmt.Printf("Could not get nodes to send: %v\n", err)
	}

	c.NodesW.Add(in)

	fmt.Println("Connected to ", in.GetId(), ":", in.GetMembershipIp(), ":", in.GetMembershipPort())
	return sendingNodes, nil
}

func (c *Connector) Disconnect(ctx context.Context, in *pb.Node) (*pb.Empty, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if err := utils.ContextError(ctx); err != nil {
		return nil, err
	}

	c.NodesW.Remove(in.GetId())

	fmt.Println("Disconnected from ", in.GetId(), ":", in.GetMembershipIp(), ":", in.GetMembershipPort())
	return nil, nil
}
