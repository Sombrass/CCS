package shadowbins

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	ShadowIdentifier string `toml:"shadow_identifier"`
	Hyper            hyperConfig
	Shadows          map[string]shadowConfig
}

type hyperConfig struct {
	Port     int
	Binaries map[string]string
}

type shadowConfig struct {
	AutoAlias []string            `toml:"auto_alias"`
	DirMap    []map[string]string `toml:"dir_map"`
}

func getConfPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(usr.HomeDir, ".shadowbins.toml")
}

func readConfig() tomlConfig {
	path := getConfPath()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	conf, err := decodeConfig(string(data))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf)
	fmt.Println(conf.Hyper.Binaries)
	fmt.Println(conf.Shadows["vagrant"].AutoAlias)
	return conf
}

func decodeConfig(data string) (tomlConfig, error) {
	var conf tomlConfig
	_, err := toml.Decode(string(data), &conf)
	return conf, err
}
