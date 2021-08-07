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

	//Open the config file
	b, err := os.ReadFile(yamlTxt)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Unmarshalling")

	//Deserialize the config file
	var conf structs.WgConfig
	err = yaml.Unmarshal(b, &conf)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return conf
}

func parseIni(conf structs.WgConfig) map[string]*ini.File {
	var iniFiles = make(map[string]*ini.File, len(conf.Nodes))

	for _, node := range conf.Nodes {

		fmt.Println("Parsing peer", node.Name)

		var iniFile = ini.Empty()

		//Add Interface section
		sec, _ := iniFile.NewSection("Interface")

		sec.NewKey("Address", node.Address)
		sec.NewKey("PrivateKey", node.PrivateKey)
		sec.NewKey("DNS", conf.Dns)

		//If the node is a lighthouse add the keys PostUp and PostDown
		if node.Lighthouse {
			sec.NewKey("PostUp", node.PostUp)
			sec.NewKey("PostDown", node.PostDown)
		}

		//if peers are directly connected to this node specify a port to listen to
		if hasConnections(node.Name, conf.Nodes) || node.Lighthouse {
			sec.NewKey("ListeningPort", node.ListeningPort)
		}

		for i, peer := range conf.Nodes {
			//A peer must be added if it's a lighthouse, if it's connected to this node or if the node is a lighthouse
			if (peer.Lighthouse || contains(node.Name, peer.ConnectedTo) || node.Lighthouse) && node.Name != peer.Name {
				sec, _ = iniFile.NewSection("Peer" + fmt.Sprint(i))
				sec.NewKey("PublicKey", peer.PublicKey)

				//if the node must connect to the peer add the endpoint and the keep alive rate
				if peer.Lighthouse || contains(node.Name, peer.ConnectedTo) {
					sec.NewKey("Endpoint", peer.Endpoint + ":" + peer.ListeningPort)
					sec.NewKey("PersistentKeepalive", fmt.Sprint(conf.PersistentKeepAlive))
				}
				//if the peer receive a direct connection from this node add the ip range to redirect, wich should exclusivly contain the address of the peer
				if node.Lighthouse || contains(node.Name, peer.ConnectedTo) {
					sec.NewKey("AllowedIps", peer.Address)
				}
				//If the peer is a light house add the defined ip range to redirect
				if peer.Lighthouse {
					sec.NewKey("AllowedIps", peer.AllowedIps)
				}
			}
		}

		iniFiles[node.Name] = iniFile
	}

	return iniFiles
}

func remvoveDigits(inis map[string]*ini.File) map[string]string {
	iniTxts := make(map[string]string, len(inis))
	//Target "Peer" plus any number of digit
	re := regexp.MustCompile("Peer[0-9]+")

	for name, iniFile := range inis {
		buffer := bytes.Buffer{}
		iniFile.WriteTo(&buffer)

		//Replace the targeted string by only "Peer"
		iniTxt := re.ReplaceAll(buffer.Bytes(), []byte("Peer"))

		iniTxts[name] = string(iniTxt)
	}

	return iniTxts
}

func contains(s string, array []string) bool {
	//Return true if an array of string contains a given string
	for _, str := range array {
		if str == s {
			return true
		}
	}
	return false
}

func hasConnections(name string, peers []structs.Node) bool {
	//return true if at least one peer has the name of the node in its ConnectTo list
	for _, peer := range peers {
		if contains(name, peer.ConnectedTo) {
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

	//Write each ini into a file with the given name of the node
	for name, iniTxt := range inisText {
		f, _ := os.Create(name + ".conf")
		f.WriteString(iniTxt)
	}
}
