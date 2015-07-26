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

		Outer:
			for _, expectedConfig := range expected {
				for _, config := range configs {
					// This is so that we don't have to worry about database assigned ids
					expectedConfig.ID = config.ID
					if config == expectedConfig {
						continue Outer
					}
				}
				return failure{Prefix: fmt.Sprintf("Config %#v not retrieved", expectedConfig)}
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
			{Name: "Config9", HostName: "Config.9", Port: 9, Username: "user9"},
		},
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
