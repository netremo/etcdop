package main

import (
	// "github.com/davecgh/go-spew/spew"
	"log"
	"time"
	"context"
	"github.com/coreos/etcd/client"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"flag"
	"os"
	// "encoding/json"
)

type Config struct{
	url string
	root_key string
	source_file string
	output_file string
	file_type string
}

type ENode struct{
	Key string	`json:"key"`
	Dir bool	`json:"dir,omitempty"`
	Value string	`json:"value,omitempty"`
	Nodes []ENode	`json:"nodes,omitempty"`	
}

func save_yaml(node ENode, filename string) error {
	node_yaml, err := yaml.Marshal(node)
	err = ioutil.WriteFile(filename, node_yaml, 0644)
	return err
}

func load_yaml(filename string) ENode {
	bytes, _ := ioutil.ReadFile(filename)
	var lnode ENode
	yaml.Unmarshal(bytes, &lnode)
	return lnode
}

func get_node(Kapi client.KeysAPI,key string) ENode {
	log.Printf("Processing key -> %q", key)
	resp, err := Kapi.Get(context.Background(), key, &client.GetOptions{Recursive:false})
	if err != nil {
		log.Fatal(err)
	} 
	node := ENode{Dir: resp.Node.Dir, Key: resp.Node.Key, Value: resp.Node.Value}
	for _ , v := range resp.Node.Nodes {
		node.Nodes = append(node.Nodes, get_node(Kapi, v.Key))
	}
	return node
}


func write_node(Kapi client.KeysAPI, node ENode) {
	log.Printf("Processing key -> %q", node.Key)
	if node.Dir {
		log.Println("Directory ", node.Key)
	} else {
		_, err := Kapi.Set(context.Background(), node.Key, node.Value, &client.SetOptions{Dir: node.Dir})
		if err != nil{
			println("Write fails on key", err.Error)
		}
	}
	
	for _, v := range node.Nodes {
		write_node(Kapi, v)
	}
}

func main() {
	urlPtr := flag.String("url", "", "ETCD url. (example http://127.0.0.1:4001)")
	rootkeyPtr := flag.String("key", "/", "Root key to use")
	source_filePtr := flag.String("in", "", "Source file to use")
	output_filePtr := flag.String("out", "", "Output file to use")
	// typePtr := flag.String("type", "yaml", "Json or yaml")
	flag.Parse()

	if (*source_filePtr=="") && (*output_filePtr=="") {
		println("ERROR: No input or output file. -h for help")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if (*urlPtr == "") {
		println("ERROR: Please specify ETCD URL. -h for help")
		flag.PrintDefaults()
		os.Exit(1)
	}
	
	work_conf := Config{
		url: *urlPtr,
		root_key: *rootkeyPtr,
		source_file: *source_filePtr,
		output_file: *output_filePtr,	
		// TODO: file_type: *typePtr
	}

	cfg := client.Config{
		Endpoints: []string{work_conf.url},
		Transport: client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := client.NewKeysAPI(c)
	
	if work_conf.source_file != "" {
		println("Starting to restore data to ETCD")
		write_node(kapi, load_yaml(work_conf.source_file))
	}

	if work_conf.output_file != "" {
		println("Starting to dump data from ETCD")
		node_to_save := get_node(kapi, work_conf.root_key)
		err := save_yaml(node_to_save, work_conf.output_file)
		if err!=nil {
			log.Fatal(err)
		}
	}
}