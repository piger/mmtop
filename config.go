package sqlgo

import (
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
)

const MYSQL_PORT int = 3306

type DbConnectionData struct {
	Name     string
	Address  string
	Port     int
	Username string
	Password string
}

func (c DbConnectionData) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)", c.Username, c.Password, c.Address, c.Port)
}

func ReadConfig(filename string) ([]DbConnectionData, error) {
	var result []DbConnectionData

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return nil, err
	}

	cfg, err := ini.Load(data)
	if err != nil {
		fmt.Printf("Error loading INI file: %s\n", err)
		return nil, err
	}

	for _, section := range cfg.Sections() {
		// fmt.Printf("section: %s\n", section.Name())
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
			fmt.Println(err)
			return nil, err
		}

		password, err := section.GetKey("password")
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		var port int = MYSQL_PORT
		if section.HasKey("port") {
			port, err = section.Key("port").Int()
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		}

		dbconn := DbConnectionData{
			Name:     name,
			Address:  hostname,
			Port:     port,
			Username: username.String(),
			Password: password.String(),
		}
		result = append(result, dbconn)
	}

	// XXX
	return result, nil
}
