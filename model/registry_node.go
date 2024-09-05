package model

import "github.com/AlessandroFinocchi/sdcc_common/pb"

type RegistryNode struct {
	Node        pb.Node
	Tick        uint64
	LastUpdated uint64
}

func (rn *RegistryNode) GetId() string {
	return rn.Node.Id
}
