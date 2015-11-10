package usage

import (
	"fmt"
	"time"

	nomadapi "github.com/hashicorp/nomad/api"
)

// Current returns the current allocation data for the datacenter
func (u *Usage) Current(datacenter string) *Snapshot {
	s := newSnapshot()
	if _, exists := u.Datacenters[datacenter]; exists == false {
		return s
	}

	for _, node := range u.Datacenters[datacenter].Nodes {
		u.node(datacenter, node, s)
	}

	return s
}

// NodeCurrent returns the allocation data for a node
func (u *Usage) NodeCurrent(datacenter, nodeName string) *Snapshot {
	s := newSnapshot()
	if _, exists := u.Datacenters[datacenter]; exists == false {
		return s
	}
	if _, exists := u.Datacenters[datacenter].Nodes[nodeName]; exists == false {
		return s
	}

	u.node(datacenter, u.Datacenters[datacenter].Nodes[nodeName], s)
	return s
}

// node fills the Snapshot struct with the requested node data
func (u *Usage) node(datacenter string, node *nomadapi.Node, s *Snapshot) {
	s.Resources.CPU += node.Resources.CPU
	s.Resources.MemoryMB += node.Resources.MemoryMB
	s.Resources.DiskMB += node.Resources.DiskMB
	s.Resources.IOPS += node.Resources.IOPS
	for _, allocation := range u.Datacenters[datacenter].Allocations[node.ID] {
		if allocation.ClientStatus == "running" {
			s.Allocated.CPU += allocation.Resources.CPU
			s.Allocated.MemoryMB += allocation.Resources.MemoryMB
			s.Allocated.DiskMB += allocation.Resources.DiskMB
			s.Allocated.IOPS += allocation.Resources.IOPS
		}
	}
}

// Loop start the loop that gets the usage information
func (u *Usage) Loop() {
	// query for the first time
	u.log.Debugln("Query Nodes Loop")
	u.queryNodes(true)
	// updates on tick since theorically nodes meta data can change too much
	// block on a query can overhead servers with node requests
	ticker := time.NewTicker(time.Duration(60) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				u.log.Debugln("Query Nodes Loop")
				u.queryNodes(true)
			case <-u.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (u *Usage) queryNodes(watch bool) {
	nodes, _, err := u.client.Nodes().List(&nomadapi.QueryOptions{})
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	err = u.UpdateNodes(nodes, watch)
	if err != nil {
		fmt.Println("err:", err)
	}
}

// UpdateNodes adds watches for new nodes
func (u *Usage) UpdateNodes(nodes []*nomadapi.NodeListStub, watch bool) error {
	for _, node := range nodes {
		// only if not exists
		if _, exists := u.Datacenters[node.Datacenter]; exists == false {
			u.Datacenters[node.Datacenter] = newDatacenter()
		}

		if node.Status == "ready" {
			if u.Datacenters[node.Datacenter].Nodes[node.ID] == nil {
				q := &nomadapi.QueryOptions{}
				data, _, err := u.client.Nodes().Info(node.ID, q)
				if err != nil {
					fmt.Println("err: ", err)
					return err
				}
				u.Datacenters[node.Datacenter].Nodes[node.ID] = data
				if watch == true {
					go u.WatchAllocationsForNode(node.Datacenter, node.ID)
				}
			}
		}
	}

	return nil
}

// WatchAllocationsForNode initializes a watch on the requested node allocations
// it stops on 404 or non-ready status
func (u *Usage) WatchAllocationsForNode(dc, id string) {
	var waitIndex uint64
	for {
		node, _, err := u.client.Nodes().Info(id, &nomadapi.QueryOptions{})
		// node is not ready? not even here?
		if node.Status != "ready" || err != nil {
			u.log.WithField("node", id).Warningln("non-ready Err:", err)
			delete(u.Datacenters[dc].Nodes, id)
			break
		}

		q := &nomadapi.QueryOptions{
			AllowStale: false,
			WaitIndex:  waitIndex,
			WaitTime:   u.waitTime,
		}

		data, meta, err := u.client.Nodes().Allocations(id, q)
		// if node goes stop monitoring it
		if err != nil {
			fmt.Println("err: ", err)
			u.log.Warningln("Allocation err")
		}

		// if index didn't change do not update anything
		if waitIndex != meta.LastIndex {
			waitIndex = meta.LastIndex
			u.Datacenters[dc].Allocations[id] = data
		}
		u.log.WithField("node", id).Debugln("updated")
		time.Sleep(10 * time.Second) // remove on blocking calls, waiting for 2.0
	}
}
