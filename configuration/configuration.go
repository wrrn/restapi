package configuration

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

const (
	uniqueViolation = "23505"
)

var DuplicateConfigErr = fmt.Errorf("Configuration exists with the same name")
var DoesNotExistErr = fmt.Errorf("Configuration does not exist")

type ConfigurationController struct {
	*sql.DB
}

type ConfigurationError struct {
	Err error
	Configuration
}

func (ce ConfigurationError) Error() string {
	return fmt.Sprintf("Configuration Error:\n Error: %s\n Configuration: %#v", ce.Err.Error(), ce.Configuration)
}

type Configuration struct {
	ID       int    `json:"id,omitempty""`
	Name     string `json:"name,omitempty"`
	HostName string `json:"hostname,omitempty"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
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

// Add attempts to add all of the configurations in the argument to the database. It returns
// a list of the names of the configurations that have been added. It return a
// ConfigurationError with an Err of DuplicateConfigError on the addition of a configuration
// that has the same name of an existing configuration.
func (cc *ConfigurationController) Add(configs ...Configuration) (names []string, err error) {
	var (
		tx   *sql.Tx
		stmt *sql.Stmt
	)
	tx, err = cc.DB.Begin()
	if err != nil {
		tx.Rollback()
		return names, err
	}

	stmt, err = tx.Prepare("INSERT INTO configurations(config_name, host_name, username, port) VALUES($1,$2,$3,$4)")
	if err != nil {
		tx.Rollback()
		return names, err
	}

	for _, config := range configs {
		_, err = stmt.Exec(config.Name, config.HostName, config.Username, config.Port)

		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolation {
				err = ConfigurationError{
					Err:           DuplicateConfigErr,
					Configuration: config,
				}
			}
			tx.Rollback()
			return names, err
		}
		names = append(names, config.Name)
	}

	err = tx.Commit()
	return names, err

}

// Delete will delete all of the configurations whose name is in the list
// of names in the arugment. It will not return an error if the name is not found.
func (cc *ConfigurationController) Delete(names ...string) (err error) {
	var (
		tx   *sql.Tx
		stmt *sql.Stmt
	)
	tx, err = cc.DB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err = tx.Prepare("DELETE FROM configurations where config_name = $1")
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, name := range names {
		_, err := stmt.Exec(name)

		if err == sql.ErrNoRows {
			err = nil
		}
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()

}

// Modify modifies the field of configuration with the same name as the name
// argument to match the fields of the second argument. All fields that are
// not set will retain their values. A configuration with the updated fields
// will be returned.
func (cc *ConfigurationController) Modify(name string, config Configuration) (newConfig Configuration, err error) {
	var (
		tx           *sql.Tx
		actualConfig Configuration
	)
	tx, err = cc.DB.Begin()
	if err != nil {
		tx.Rollback()
		return newConfig, err
	}

	err = tx.QueryRow("SELECT id, config_name, host_name, username, port FROM configurations WHERE config_name = $1", name).Scan(&actualConfig.ID, &actualConfig.Name, &actualConfig.HostName, &actualConfig.Username, &actualConfig.Port)
	if err == sql.ErrNoRows {
		err = DoesNotExistErr
	}
	if err != nil {
		tx.Rollback()
		return newConfig, err
	}

	if config.Name == "" {
		config.Name = actualConfig.Name
	}

	if config.HostName == "" {
		config.HostName = actualConfig.HostName
	}

	if config.Username == "" {
		config.Username = actualConfig.Username
	}

	if config.Port == 0 {
		config.Port = actualConfig.Port
	}

	_, err = tx.Exec(
		`UPDATE configurations 
         SET 
           config_name = $1,
           host_name = $2,
           username = $3,
           port = $4
        WHERE
          config_name = $5`, config.Name, config.HostName, config.Username, config.Port, name)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolation {
			err = ConfigurationError{
				Err:           DuplicateConfigErr,
				Configuration: config,
			}
		} else if err == sql.ErrNoRows {
			err = DoesNotExistErr
		}
		tx.Rollback()
		return newConfig, err
	}

	tx.Commit()
	return config, nil

}

func EqualConfigurations(x, y Configuration) bool {
	return x.Name == y.Name &&
		x.HostName == y.HostName &&
		x.Username == y.Username &&
		x.Port == y.Port
}
