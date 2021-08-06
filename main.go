package main

import (
	"fmt"
	"os"

	// ini "gopkg.in/ini.v1"
	"github.com/LeoDPlouc/wg-autoconfig/structs"
	yaml "gopkg.in/yaml.v2"
)

func parseYaml(yamlTxt string) structs.WgConfig {
	b, err := os.ReadFile(yamlTxt)

	if err == nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var conf structs.WgConfig
	yaml.Unmarshal(b, conf)

	return conf
}

func parseIni(conf structs.WgConfig) structs.IniFile {
}

func main() {
	var yamlFile = os.Args[1]
	/*var wgConf = */ parseYaml(yamlFile)
	//var ini = parseIni(wgConf)
}
