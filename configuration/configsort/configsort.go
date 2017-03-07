package configsort

import (
	"sort"

	"github.com/wrrn/restapi/configuration"
)

type Orderer func(i, j int) bool

// Sorter sorts functions based on it's Order
type Sorter struct {
	Configs []configuration.Configuration
	Order   Orderer
}

// Sort sorts the configs according to the order that was given.
func (s *Sorter) Sort(o Orderer, configs []configuration.Configuration) []configuration.Configuration {
	s.Configs = configs
	s.Order = o
	sort.Sort(s)
	return s.Configs
}

// Sort reverse sorts the configs according to the order that was given.
func (s *Sorter) Reverse(o Orderer, configs []configuration.Configuration) []configuration.Configuration {
	s.Configs = configs
	s.Order = o
	sort.Reverse(s)
	return s.Configs
}

func (s Sorter) Len() int {
	return len(s.Configs)
}

func (s Sorter) Swap(i, j int) {
	temp := s.Configs[i]
	s.Configs[i] = s.Configs[j]
	s.Configs[j] = temp
}

func (s Sorter) Less(i, j int) bool {
	return s.Order(i, j)
}

// ByName returns true if the name of the configuration at index i is less than
// the name of the configuration at index j.
func (s *Sorter) ByName(i, j int) bool {
	return s.Configs[i].Name < s.Configs[j].Name
}

// ByHostName returns true if the HostName of the configuration at index i is less than
// the HostName of the configuration at index j.
func (s *Sorter) ByHostName(i, j int) bool {
	return s.Configs[i].HostName < s.Configs[j].HostName
}

// ByPort returns true if the port of the configuration at index i is less than
// the port of the configuration at index j.
func (s *Sorter) ByPort(i, j int) bool {
	return s.Configs[i].Port < s.Configs[j].Port
}

// ByUsername returns true if the username of the configuration at index i is less than
// the username of the configuration at index j.
func (s *Sorter) ByUsername(i, j int) bool {
	return s.Configs[i].Username < s.Configs[i].Username
}
