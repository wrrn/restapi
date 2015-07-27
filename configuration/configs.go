package configuration

import "sort"

type Orderer func(i, j int) bool

type Configurations struct {
	Configs []Configuration `json:"configurations"`
}

// GetFirst retrieves the first configuration from the Configurations object.
// If the length Configurations object's Configs is less than 1 it returns an
// uninitialized struct.
func (cs Configurations) GetFirst() (config Configuration) {
	if len(cs.Configs) > 0 {
		return cs.Configs[0]
	}
	return config
}

// Sorter sorts functions based on it's Order
type Sorter struct {
	Configs []Configuration
	Order   Orderer
}

// Equals determines if two Configuration slices are equal
func Equals(x, y []Configuration) bool {
	if len(x) != len(y) {
		return false
	}

	for _, configX := range x {
		var found bool
		for _, configY := range y {

			if EqualConfigurations(configX, configY) {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}
	return true
}

// Sort sorts the configs according to the order that was given.
func (s *Sorter) Sort(o Orderer, configs []Configuration) []Configuration {
	s.Order = o
	sort.Sort(s)
	return s.Configs
}

// Sort reverse sorts the configs according to the order that was given.
func (s *Sorter) Reverse(o Orderer, configs []Configuration) []Configuration {
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
func (s Sorter) ByName(i, j int) bool {
	return s.Configs[i].Name < s.Configs[j].Name
}

// ByHostName returns true if the HostName of the configuration at index i is less than
// the HostName of the configuration at index j.
func (s Sorter) ByHostName(i, j int) bool {
	return s.Configs[i].HostName < s.Configs[j].HostName
}

// ByPort returns true if the port of the configuration at index i is less than
// the port of the configuration at index j.
func (s Sorter) ByPort(i, j int) bool {
	return s.Configs[i].Port < s.Configs[j].Port
}

// ByUsername returns true if the username of the configuration at index i is less than
// the username of the configuration at index j.
func (s Sorter) ByUsername(i, j int) bool {
	return s.Configs[i].Username < s.Configs[i].Username
}
