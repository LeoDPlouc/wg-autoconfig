# wg-autoconfig

Here is a tool to write wireguard interface for all your nodes with a config file.

## How to use it ?
Write your config file (more on that later).
<br>Run `wg-config <your-config-file>`

## How to write the config file ?

The config file is written in yaml.
<br>The fields of the config file are the same as wireguard config file.
<br>Some fields are specific to this tool

- Name : This field is used to differentiate the nodes in a user friendly way
- Lighthouse : Must be a bool. A light house is a node with a static IP accessible via internet. All the nodes will be connected with all the lighthouses. This assure that all nodes are accessible by all the others node as long as all the lighthouses are alive.
- ConnectedTo : If some nodes are directly accessible to other nodes they can be peered. For example two nodes living in the same local network

#### /!\ ALL KEYS MUST BE WRITTEN LOWERCASE /!\


Here is an example with a VPS accessible via internet with a static IP, a raspberry and a desktop living in the same local network with static IP and a smartphone with a dynamic IP and network.

```yaml
dns: 1.1.1.1
persistentkeepalive: 50
peers:
  - name: Server
    publickey: ExamplePubKey1
    privatekey: ExamplePriKey1
    address: 10.0.0.1/32
    endpoint: 10.1.1.1
    lighthouse: true
    allowedips: 10.0.0.0/24
    postup: iptables -A FORWARD -i %i -j ACCEPT; iptables -A FORWARD -o %i -j ACCEPT; iptables -t nat -A POSTROUTING -o ens3 -j MASQUERADE
    postdown: iptables -D FORWARD -i %i -j ACCEPT; iptables -D FORWARD -o %i -j ACCEPT; iptables -t nat -D POSTROUTING -o ens3 -j MASQUERADE

  - name: Smartphone
    publickey: ExamplePubKey2
    privatekey: ExamplePriKey2
    address: 10.0.0.2/32
    allowedips: 10.0.0.0/24

  - name: Desktop
    publickey: ExamplePubKey3
    privatekey: ExamplePriKey3
    address: 10.0.0.3/32
    endpoint: 10.2.2.2
    connectedto:
      - RaspberryPi
    allowedips: 10.0.0.0/24

  - name: RaspberryPi
    publickey: ExamplePubKey4
    privatekey: ExamplePriKey4
    address: 10.0.0.4/32
    endpoint: 10.3.3.3
    connectedto:
      - Desktop
    allowedips: 10.0.0.0/24

```

## How to build it ?

This tool is a standard go program and must be compile like every go program.

Make sure go is installed 
<br>Clone this repo under $GOPATH/src
<br>Open a terminal in the cloned repo
> go build