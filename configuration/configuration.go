package configuration

import "database/sql"

type ConfigurationController struct {
	*sql.DB
}

type Configuration struct {
	ID       int
	Name     string
	HostName string
	Port     int
	Username string
}

// GetAll returns a list of all of the stored configurations
func (cc *ConfigurationController) GetAll() (configs []Configuration, err error) {
	rows, err := cc.DB.Query("SELECT id, config_name, host_name, username, port FROM configurations")
	if err == sql.ErrNoRows {
		return configs, nil
	} else if err != nil {
		return configs, err
	}

	for rows.Next() {
		config := Configuration{}
		err = rows.Scan(&config.ID, &config.Name, &config.HostName, &config.Username, &config.Port)
		if err == nil {
			configs = append(configs, config)
		}
	}
	return configs, rows.Err()
}
