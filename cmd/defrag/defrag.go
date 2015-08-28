package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"

	"github.com/Patrolavia/ggpk/afs"
	"github.com/Patrolavia/ggpk/record"
)

func W(f *os.File, data interface{}) error {
	return binary.Write(f, binary.LittleEndian, data)
}

func main() {
	flag.Parse()
	fn := flag.Arg(0)

	done := func(err error) {
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		fmt.Print("done.\n")
	}

	fmt.Printf("Loading GGPK file from %s ... ", fn)
	orig, err := os.Open(fn)
	done(err)

	fmt.Print("Reading GGPK content ... ")
	root, err := afs.FromGGPK(orig)
	done(err)

	fmt.Print("Creating result.ggpk")
	dest, err := os.Create("result.ggpk")
	defer dest.Close()
	done(err)

	fmt.Print("Writing signature ... ")
	// write sign
	ggg := record.GGGRecord{}
	ggg.Offsets = make([]uint64, 1)
	ggg.NodeCount = 1
	ggg.Header.Length = uint32(ggg.ByteLength())
	ggg.Header.Tag = "GGPK"
	ggg.Offsets[0] = uint64(ggg.ByteLength())
	done(ggg.Save(dest))


	dirs, files := compute(root, uint64(ggg.Header.Length))
	_ = dirs
	_ = files

	size := ggg.Header.Length
	for _, d := range dirs {
		size += d.h.Length
	}

	totalDirCount := 0
	for range dirs {
		totalDirCount++
	}
	totalFileBytes := uint64(0)
	for _, f := range files {
		totalFileBytes += uint64(f.h.Length)
	}

	level := 0
	p := func(idx, max uint64) {
		progress := idx * 10 / max
		if level < int(progress) {
			level = int(progress)
			fmt.Printf("%d0%% ", level)
		}
	}

	fmt.Print("Writing directory structures ... ")
	level = 0
	for idx, d := range dirs {
		d.save(dest)
		p(uint64(idx), uint64(totalDirCount))
	}
	done(nil)

	fmt.Print("Writing file contents ... ")
	level = 0
	cur := uint64(0)
	for _, f := range files {
		f.save(dest, orig)
		cur += uint64(f.h.Length)
		p(cur, totalFileBytes)
	}
	done(nil)
}
