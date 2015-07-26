package configuration

type Configs []Configuration

func Equals(x []Configuration, y []Configuration) bool {
	if len(x) != len(y) {
		return false
	}

	for _, configX := range x {
		var found bool
		for _, configY := range y {

			// This is so that we don't have to worry about database assigned ids

			if configX.Name == configY.Name &&
				configX.HostName == configY.HostName &&
				configX.Username == configY.Username &&
				configX.Port == configY.Port {
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
