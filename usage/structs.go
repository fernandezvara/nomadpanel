package usage

import (
	"time"

	"github.com/Sirupsen/logrus"
	nomadapi "github.com/hashicorp/nomad/api"
)

// Usage has the nodes information on memory and updates it using watches
// It's updated when the index changes, not by pulling all data each time
type Usage struct {
	client      *nomadapi.Client
	Datacenters map[string]*Datacenter
	quit        chan struct{}
	waitTime    time.Duration
	log         *logrus.Logger
}

// NewUsage returns a usage client
func NewUsage(client *nomadapi.Client, waitTime time.Duration, log *logrus.Logger) *Usage {
	return &Usage{
		client:      client,
		Datacenters: make(map[string]*Datacenter),
		waitTime:    waitTime,
		log:         log,
	}
}

func newDatacenter() *Datacenter {
	return &Datacenter{
		Nodes:       make(map[string]*nomadapi.Node),
		Allocations: make(map[string][]*nomadapi.Allocation),
	}
}

// Datacenter contains the information to match nodes/allocs
type Datacenter struct {
	Nodes       map[string]*nomadapi.Node
	Allocations map[string][]*nomadapi.Allocation
}

func newSnapshot() *Snapshot {
	return &Snapshot{
		Resources: new(Resources),
		Allocated: new(Resources),
	}
}

// Snapshot is the struct that returns the current usage
type Snapshot struct {
	Resources *Resources `json:"resources"`
	Allocated *Resources `json:"allocated"`
}

// Resources is the nomad exposed resources per node
type Resources struct {
	CPU      int
	MemoryMB int `mapstructure:"memory"`
	DiskMB   int `mapstructure:"disk"`
	IOPS     int
}
