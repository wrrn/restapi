package configuration

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

// GetPage will return a Configuration slice of the Configurations that would
// be on the defined page. Note: pageNum are 0 indexed.
func GetPage(configs []Configuration, pageNum, perPage int) []Configuration {
	start := pageNum * perPage
	end := start + perPage
	if start > len(configs) || pageNum < 0 || perPage < 0 {
		return make([]Configuration, 0)
	}

	if end > len(configs) {
		end = len(configs)
	}
	configsPage := make([]Configuration, end-start)
	copy(configsPage, configs[start:end])
	return configsPage
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
