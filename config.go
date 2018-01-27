package mmtop

import (
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
)

const MYSQL_PORT int = 3306

type DbConnectionInfo struct {
	Name     string
	Address  string
	Port     int
	Username string
	Password string
}

func (c DbConnectionInfo) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)", c.Username, c.Password, c.Address, c.Port)
}

func ReadConfig(filename string) ([]DbConnectionInfo, error) {
	var configs []DbConnectionInfo

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg, err := ini.Load(data)
	if err != nil {
		return nil, err
	}

	for _, section := range cfg.Sections() {
		var hostname string

		name := section.Name()
		if name == "DEFAULT" {
			continue
		}

		if section.HasKey("hostname") {
			hostname = section.Key("hostname").String()
		} else {
			hostname = name
		}

		username, err := section.GetKey("username")
		if err != nil {
			return nil, err
		}

		password, err := section.GetKey("password")
		if err != nil {
			return nil, err
		}

		var port int = MYSQL_PORT
		if section.HasKey("port") {
			port, err = section.Key("port").Int()
			if err != nil {
				return nil, err
			}
		}

		dbconn := DbConnectionInfo{
			Name:     name,
			Address:  hostname,
			Port:     port,
			Username: username.String(),
			Password: password.String(),
		}
		configs = append(configs, dbconn)
	}

	return configs, nil
}
