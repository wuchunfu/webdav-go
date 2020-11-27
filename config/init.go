package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type DavServer struct {
	Ip   string `yaml:"ip,omitempty"`
	Port uint16 `yaml:"port"`

	Auth bool   `yaml:"auth"`
	User string `yaml:"user,omitempty"`
	Pass string `yaml:"pass,omitempty"`

	Scope  string `yaml:"scope"`
	Modify bool   `yaml:"modify"`

	Tls  bool   `yaml:"tls"`
	Cert string `yaml:"cert,omitempty"`
	Key  string `yaml:"key,omitempty"`
}
type MainConfig struct {
	Log     string      `yaml:"log"`
	Default DavServer   `yaml:"default"`
	Servers []DavServer `yaml:"server"`
}

var GlobalConf MainConfig

func ReloadConfig() {
	GlobalConf = MainConfig{
		Log:     "webdav.log",
		Default: DavServer{},
		Servers: []DavServer{},
	}

	config, err := ioutil.ReadFile("webdav.yaml")
	if err != nil {
		log.Fatalln("No config file.")
		return
	}

	err = yaml.Unmarshal(config, &GlobalConf)
	if err != nil {
		log.Println("Please check config.")
		log.Fatalln(err.Error())
		return
	}

	fd, err := os.OpenFile(GlobalConf.Log, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Cannot open %s.\n", GlobalConf.Log)
		return
	}
	log.SetOutput(fd)

	GlobalConf.Default.Auth = true
	GlobalConf.Default.Tls = true
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	ReloadConfig()
}
