package main

import (
	"ggpk/record"
	"log"
	"os"
)

func main() {
	f, err := os.Open("Content.ggpk")
	if err != nil {
		log.Fatalf("Cannot open Content.ggpk: %s", err)
	}
	defer f.Close()

	root, err := record.GGG(f)
	if err != nil {
		log.Fatalf("Cannot read data from Content.ggpk: %s", err)
	}

	if string([]byte(root.Header.Tag)) != "GGPK" {
		log.Fatalf("Expected GGPK, got %s", root.Header.Tag)
	}

	childnodes, err := root.Children(f)
	if err != nil {
		log.Fatalf("Cannot read root nodes from Content.ggpk: %s", err)
	}
	log.Printf("%#v", childnodes)
}
