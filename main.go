package main

import (
	"fmt"
	"os"

	// ini "gopkg.in/ini.v1"
	yaml "gopkg.in/yaml.v2"
	"github.com/LeoDPlouc/wg-autoconfig/structs"
	//"structs/WgConfig"
)

func parseYaml(yamlTxt string) wgConfig {
	b, err := os.ReadFile(yamlTxt)

	if err == nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var conf wgConfig
	yaml.Unmarshal(b, conf)

	return conf
}

func parseIni(conf wgConfig) iniFile {

}

func main() {
	var yamlFile = os.Args[1]
	/*var wgConf = */ parseYaml(yamlFile)
	//var ini = parseIni(wgConf)
}
