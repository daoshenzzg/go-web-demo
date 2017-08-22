package config

import (
	"github.com/go-ini/ini"
)

var (
	iniFile       *ini.File
	ServerPort    string
	ServerTimeout int
)

func init() {
	file, err := ini.Load("conf/conf.ini")
	if err != nil {
		panic(err)
	}

	iniFile = file

	// Server Port
	ServerPort = Section("http").Key("port").String()

	// Server Timeout
	ServerTimeout, _ = Section("http").Key("timeout").Int()
}

func Section(name string) *ini.Section {
	return iniFile.Section(name)
}
