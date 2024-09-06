package model

import (
	"github.com/AlessandroFinocchi/sdcc_common/pb"
)

type NodeListWrapper struct {
	NodeList []*pb.Node
	NodeMap  map[string]*RegistryNode
}

func NewNodeListWrapper() *NodeListWrapper {
	nodeList := make([]*pb.Node, 0)

	nodeMap := make(map[string]*RegistryNode)

	return &NodeListWrapper{
		NodeList: nodeList,
		NodeMap:  nodeMap,
	}
}

func (n *NodeListWrapper) Remove(id string) {
	// Remove node from the map
	delete(n.NodeMap, id)

	// Remove node from the list
	for i := range n.NodeList {
		if n.NodeList[i].Id == id {
			n.NodeList = append(n.NodeList[:i], n.NodeList[i+1:]...)
			break
		}
	}
}

func (n *NodeListWrapper) Add(node *pb.Node) {
	// Add node to the map
	n.NodeMap[node.Id] = &RegistryNode{
		Node: pb.Node{
			Id:             node.GetId(),
			MembershipIp:   node.GetMembershipIp(),
			MembershipPort: node.GetMembershipPort(),
			VivaldiIp:      node.GetVivaldiIp(),
			VivaldiPort:    node.GetVivaldiPort(),
			GossipIp:       node.GetGossipIp(),
			GossipPort:     node.GetGossipPort(),
		},
		Tick:        0,
		LastUpdated: 0,
	}

	// Add node to the list
	n.NodeList = append(n.NodeList, &pb.Node{
		Id:             node.GetId(),
		MembershipIp:   node.GetMembershipIp(),
		MembershipPort: node.GetMembershipPort(),
		VivaldiIp:      node.GetVivaldiIp(),
		VivaldiPort:    node.GetVivaldiPort(),
		GossipIp:       node.GetGossipIp(),
		GossipPort:     node.GetGossipPort(),
	})
}
