package configuration

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

type failure struct {
	Prefix   string
	Expected interface{}
	Actual   interface{}
}

func SetupDB() *sql.DB {
	db, err := sql.Open("postgres", "user=tenable password=insecure dbname=apitest")
	if err != nil {
		log.Fatal(err)
	}
	ResetDB(db)
	return db
}

func ResetDB(db *sql.DB) {
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM configurations")
	db.Exec("DELETE FROM sessions")
}

func (f failure) Error() string {
	str := f.Prefix
	if f.Expected != nil {
		str += fmt.Sprintf("\n Expected: %v", f.Expected)

	}

	if f.Actual != nil {
		str += fmt.Sprintf("\n Actual: %v", f.Actual)
	}

	return str

}

var baseExpected = []Configuration{
	{Name: "Config1", HostName: "Config.1", Port: 1, Username: "user1"},
	{Name: "Config2", HostName: "Config.2", Port: 2, Username: "user2"},
	{Name: "Config3", HostName: "Config.3", Port: 3, Username: "user3"},
	{Name: "Config4", HostName: "Config.4", Port: 4, Username: "user4"},
	{Name: "Config5", HostName: "Config.5", Port: 5, Username: "user5"},
	{Name: "Config6", HostName: "Config.6", Port: 6, Username: "user6"},
	{Name: "Config7", HostName: "Config.7", Port: 7, Username: "user7"},
	{Name: "Config8", HostName: "Config.8", Port: 8, Username: "user8"},
	{Name: "Config9", HostName: "Config.9", Port: 9, Username: "user9"},
}

var tests = map[string]struct {
	test     func(*ConfigurationController, []Configuration) error
	expected []Configuration
}{
	"TestEmptyConfigs": {
		test: func(cc *ConfigurationController, expected []Configuration) error {
			configs, err := cc.GetAll()
			if len(configs) != len(expected) {
				return failure{"Configs length does not match", len(expected), len(configs)}
			}
			return err
		},
		expected: []Configuration{},
	},

	"TestGetAll": {
		test: func(cc *ConfigurationController, expected []Configuration) error {
			for _, config := range expected {
				_, err := cc.DB.Exec(`
	                         INSERT INTO configurations(config_name, host_name, port, username) VALUES ($1, $2, $3, $4)`,
					config.Name, config.HostName, config.Port, config.Username)
				if err != nil {
					return err
				}

			}

			configs, err := cc.GetAll()
			if err != nil {
				return err
			}

			if len(configs) != len(expected) {
				return failure{"Configs length does not match", len(expected), len(configs)}
			}

			if !Equals(configs, expected) {
				return failure{Expected: expected, Actual: configs}
			}

			return nil
		},
		expected: baseExpected,
	},

	"TestAddOne": {
		test: func(cc *ConfigurationController, expected []Configuration) error {
			configs, err := cc.Add(expected...)
			names := make([]string, 0, len(expected))
			if err != nil {
				return err
			}
			for index := range configs {
				if configs[index].Name != expected[index].Name {
					return failure{"Names do not match", expected[index].Name, names[index]}
				}
				names = append(names, expected[index].Name)
			}

			configs, err = cc.Get(names...)
			if err != nil {
				return err
			}

			if !Equals(configs, expected) {
				return failure{Expected: expected, Actual: configs}
			}

			return nil
		},
		expected: baseExpected[:1],
	},

	"TestAddMultiple": {
		test: func(cc *ConfigurationController, expected []Configuration) error {
			configs, err := cc.Add(expected...)
			names := make([]string, 0, len(expected))
			if err != nil {
				return err
			}
			for index := range configs {
				if configs[index].Name != expected[index].Name {
					return failure{"Names do not match", expected[index].Name, names[index]}
				}
				names = append(names, expected[index].Name)
			}

			configs, err = cc.Get(names...)
			if err != nil {
				return err
			}

			if !Equals(configs, expected) {
				return failure{Expected: expected, Actual: configs}
			}

			return nil
		},
		expected: baseExpected,
	},
	"TestAddCollision": {
		test: func(cc *ConfigurationController, data []Configuration) error {
			_, err := cc.Add(data...)
			if err, ok := err.(ConfigurationError); !ok || err.Err != DuplicateConfigErr {
				return failure{"Errors do not match",
					ConfigurationError{
						DuplicateConfigErr,
						data[8],
					},
					err}
			}
			count := -1
			err = cc.DB.QueryRow("SELECT COUNT(id) from configurations").Scan(&count)
			if err != nil {
				return err
			}
			if count != 0 {
				return failure{"Too many configurations in DB", 0, count}
			}

			return nil
		},
		expected: []Configuration{
			{Name: "Config1", HostName: "Config.1", Port: 1, Username: "user1"},
			{Name: "Config2", HostName: "Config.2", Port: 2, Username: "user2"},
			{Name: "Config3", HostName: "Config.3", Port: 3, Username: "user3"},
			{Name: "Config4", HostName: "Config.4", Port: 4, Username: "user4"},
			{Name: "Config5", HostName: "Config.5", Port: 5, Username: "user5"},
			{Name: "Config6", HostName: "Config.6", Port: 6, Username: "user6"},
			{Name: "Config7", HostName: "Config.7", Port: 7, Username: "user7"},
			{Name: "Config8", HostName: "Config.8", Port: 8, Username: "user8"},
			{Name: "Config1", HostName: "Config.9", Port: 9, Username: "user9"},
		},
	},
	"Delete": {
		test: func(cc *ConfigurationController, data []Configuration) error {
			expected := make([]Configuration, 0, len(data))
			toDelete := make([]Configuration, 0, 4)
			expected = append(append(expected, data[0:4]...), data[7])
			toDelete = append(append(toDelete, data[4:7]...), data[8:]...)
			namesToDelete := make([]string, 0, len(toDelete))
			for _, config := range toDelete {
				namesToDelete = append(namesToDelete, config.Name)
			}

			if _, err := cc.Add(data...); err != nil {
				return err
			}

			if err := cc.Delete(namesToDelete...); err != nil {
				return err
			}

			actual, err := cc.GetAll()
			if err != nil {
				return err
			}

			if !Equals(actual, expected) {
				return failure{Expected: expected, Actual: actual}
			}

			return nil
		},
		expected: baseExpected,
	},
	"DeleteNonexisting": {
		test: func(cc *ConfigurationController, data []Configuration) error {
			if _, err := cc.Add(data...); err != nil {
				return err
			}
			if err := cc.Delete("THIS DOES NOT EXIST"); err != nil {
				return failure{"Unexpected Error", nil, err}
			}
			return nil
		},
		expected: baseExpected,
	},

	"TestModify": {
		test: func(cc *ConfigurationController, data []Configuration) error {

			_, err := cc.Add(data...)
			if err != nil {
				return err
			}

			expectedConfig := data[0]
			expectedConfig.Name = "Hello"
			newConfig := Configuration{Name: "Hello"}

			if newConfig, err = cc.Modify(data[0].Name, newConfig); err != nil {
				return err
			}

			if !EqualConfigurations(expectedConfig, newConfig) {
				return failure{Expected: expectedConfig, Actual: newConfig}
			}

			newConfigs, err := cc.Get(expectedConfig.Name)
			if err != nil {
				return err
			}

			if len(newConfigs) != 1 || !EqualConfigurations(expectedConfig, newConfigs[0]) {
				return failure{Expected: expectedConfig, Actual: newConfigs}
			}

			return nil
		},
		expected: baseExpected[:1],
	},

	"TestModifyAllFields": {
		test: func(cc *ConfigurationController, data []Configuration) error {
			_, err := cc.Add(data...)
			expectedConfig := Configuration{
				Name:     "Something else",
				HostName: "Other.stuff",
				Username: "NewUserName",
				Port:     9090,
			}
			_ = expectedConfig

			if err != nil {
				return err
			}
			newConfig := expectedConfig
			if newConfig, err = cc.Modify(data[0].Name, newConfig); err != nil {
				return err
			}

			if !EqualConfigurations(expectedConfig, newConfig) {
				return failure{Expected: data[0], Actual: newConfig}
			}

			newConfigs, err := cc.Get(expectedConfig.Name)
			if err != nil {
				return err
			}

			if len(newConfigs) != 1 || !EqualConfigurations(expectedConfig, newConfigs[0]) {
				return failure{Expected: expectedConfig, Actual: newConfigs}
			}

			return nil
		},
		expected: baseExpected[:1],
	},
	"TestModifyNonExisting": {
		test: func(cc *ConfigurationController, data []Configuration) error {
			_, err := cc.Add(data...)
			expectedConfig := Configuration{
				Name:     "Something else",
				HostName: "Other.stuff",
				Username: "NewUserName",
				Port:     9090,
			}
			_ = expectedConfig

			if err != nil {
				return err
			}
			newConfig := expectedConfig
			if newConfig, err = cc.Modify("NOT A REAL NAME", newConfig); err != DoesNotExistErr {
				return failure{Prefix: "Modify modified non-existing config", Expected: DoesNotExistErr, Actual: err}
			}

			if newConfigs, err := cc.Get(expectedConfig.Name); err != DoesNotExistErr || len(newConfigs) != 0 {
				return failure{Prefix: "Config should not have been found:", Expected: DoesNotExistErr, Actual: err}
			}

			return nil
		},
		expected: baseExpected[:1],
	},
}

func TestConfiguration(t *testing.T) {
	cc := &ConfigurationController{SetupDB()}
	for name, test := range tests {
		if err := test.test(cc, test.expected); err != nil {
			t.Errorf("%s Failed: %s", name, err.Error())
		}
		ResetDB(cc.DB)
	}

}
