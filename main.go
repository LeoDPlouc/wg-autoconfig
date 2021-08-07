package main

import (
	"fmt"
	"os"

	ini "gopkg.in/ini.v1"
	yaml "gopkg.in/yaml.v2"

	"github.com/LeoDPlouc/wg-autoconfig/structs"
)

func parseYaml(yamlTxt string) structs.WgConfig {
	fmt.Println("Opening", yamlTxt)

	b, err := os.ReadFile(yamlTxt)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Unmarshalling")

	var conf structs.WgConfig
	err = yaml.Unmarshal(b, &conf)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return conf
}

func parseIni(conf structs.WgConfig) map[string]*ini.File {
	var iniFiles = make(map[string]*ini.File, len(conf.Peers))

	for _, peer := range conf.Peers {

		fmt.Println("Parsing peer", peer.Name)

		var iniFile = ini.Empty()

		sec, _ := iniFile.NewSection("Interface")

		sec.NewKey("Address", peer.Address)
		sec.NewKey("PrivateKey", peer.PrivateKey)
		sec.NewKey("DNS", conf.Dns)

		for _, connection := range conf.Peers {
			sec, _ = iniFile.NewSection("Peer")
			
			if connection.Lighthouse || contains(peer.Name, connection.ConnectedTo) {
				sec.NewKey("PublicKey", connection.PublicKey)
				sec.NewKey("AllowedIps", connection.AllowedIps)
				sec.NewKey("Endpoint", connection.Endpoint)
				sec.NewKey("PersistentKeepalive", fmt.Sprint(conf.PersistentKeepAlive))
			}
		}

		iniFiles[peer.Name] = iniFile
	}

	return iniFiles
}

func contains(s string, array []string) bool {
	for _, str := range array {
		if str == s {
			return true
		}
	}
	return false
}

func main() {
	var yamlFile = os.Args[1]
	var wgConf = parseYaml(yamlFile)
	var inis = parseIni(wgConf)

	for name, iniFile := range inis {
		err := iniFile.SaveTo(name + ".conf")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
