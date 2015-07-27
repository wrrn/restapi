package configuration

type Configs []Configuration

func Equals(x, y []Configuration) bool {
	if len(x) != len(y) {
		return false
	}

	for _, configX := range x {
		var found bool
		for _, configY := range y {

			// This is so that we don't have to worry about database assigned ids

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
