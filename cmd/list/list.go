package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Patrolavia/ggpk/afs"
)

func main() {
	flag.Parse()
	fn := flag.Arg(0)
	f, err := os.Open(fn)
	if err != nil {
		log.Fatalf("Cannot open Content.ggpk at %s: %s", fn, err)
	}
	defer f.Close()

	root, err := afs.FromGGPK(f)
	if err != nil {
		log.Fatalf("Parse error: %s", err)
	}

	traverse(root)
}

func traverse(dir *afs.Directory) {
	fmt.Println(dir.Path)
	for _, e := range dir.Files {
		fmt.Println(e.Path)
	}

	for _, e := range dir.Subfolders {
		traverse(e)
	}
}
