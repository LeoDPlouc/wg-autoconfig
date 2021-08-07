package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

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

		for i, connection := range conf.Peers {
			if (connection.Lighthouse || contains(peer.Name, connection.ConnectedTo) || peer.Lighthouse) && peer.Name != connection.Name {
				sec, _ = iniFile.NewSection("Peer" + fmt.Sprint(i))

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

func remvoveDigits(inis map[string]*ini.File) map[string]string {
	iniTxts := make(map[string]string, len(inis))
	re := regexp.MustCompile("Peer[0-9]+")

	for name, iniFile := range inis {
		buffer := bytes.Buffer{}
		iniFile.WriteTo(&buffer)

		iniTxt := re.ReplaceAll(buffer.Bytes(), []byte("Peer"))

		iniTxts[name] = string(iniTxt)
	}

	return iniTxts
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
	var inisText = remvoveDigits(inis)

	for name, iniTxt := range inisText {
		f, _ := os.Create(name + ".conf")
		f.WriteString(iniTxt)
	}
}
