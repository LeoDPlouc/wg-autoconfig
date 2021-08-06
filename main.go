package main

import (
	"fmt"
	"os"
	// ini "gopkg.in/ini.v1"
	// yaml "gopkg.in/yaml.v2"
)

type wgConfig struct {
}

type iniFile struct {
}

func parseYaml(yamlTxt string) wgConfig {
	os.ReadFile(yamlTxt)
}

func parseIni(conf wgConfig) iniFile {

}

func main() {
	var yamlFile = os.Args[1]
	/*var wgConf = */ parseYaml(yamlFile)
	//var ini = parseIni(wgConf)
}
