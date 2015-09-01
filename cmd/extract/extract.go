package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Patrolavia/ggpk/afs"
)

var (
	recursive bool
	destDir   string
)

func init() {
	flag.BoolVar(&recursive, "r", false, "Recursive extract directory, ignored if extracting file.")
	flag.StringVar(&destDir, "d", ".", "Extract files to directory `N`.")
	flag.Parse()
}

func main() {
	fn := flag.Arg(0)
	path := flag.Arg(1)
	if path == "" {
		log.Fatalf("You have to specify path to extract.")
	}
	if path[len(path)-1] == "/"[0] {
		path = path[:len(path)-1]
	}
	f, err := os.Open(fn)
	if err != nil {
		log.Fatalf("Cannot open ggpk file at %s: %s", fn, err)
	}
	defer f.Close()

	log.Print("Parsing ggpk file ...")
	root, err := afs.FromGGPK(f)
	if err != nil {
		log.Fatalf("Parse error: %s", err)
	}

	if path == "" {
		saveDir(root, f)
		return
	}

	cur := root
	nodes := strings.Split(path[1:], "/")
Orz:
	for idx, node := range nodes {
		if idx == len(nodes)-1 {
			for _, file := range cur.Files {
				if file.Name == node {
					saveFile(file, f)
					return
				}
			}
		}

		for _, dir := range cur.Subfolders {
			if dir.Name == node {
				if idx == len(nodes)-1 {
					saveDir(dir, f)
					return
				}
				cur = dir
				continue Orz
			}
		}
		log.Fatalf("Cannot find %s in %s", node, cur.Path)

	}
}

func saveFile(file *afs.File, f *os.File) {
	fmt.Printf("Writing file %s ... ", file.Path)
	data, err := file.Content()
	if err != nil {
		log.Fatalf("While reading file %s: %s", file.Path, err)
	}

	fn := filepath.FromSlash(destDir + file.Path)
	dirname := filepath.Dir(fn)
	if err := os.MkdirAll(dirname, os.FileMode(0777)); err != nil {
		log.Fatalf("Cannot create directory %s: %s", dirname, err)
	}

	dest, err := os.Create(fn)
	if err != nil {
		log.Fatalf("Error creating file %s: %s", file.Path, err)
	}
	defer dest.Close()

	if _, err := dest.Write(data); err != nil {
		log.Fatalf("Error writing file %s: %s", file.Path, err)
	}
	fmt.Printf("%d bytes\n", file.Size)
}

func saveDir(dir *afs.Directory, f *os.File) {
	for _, file := range dir.Files {
		saveFile(file, f)
	}

	if recursive {
		for _, child := range dir.Subfolders {
			saveDir(child, f)
		}
	}
}
