package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/marcinbor85/gohex"
)

const defaultPadByte = 0x00

func main() {
	var pad uint

	flag.UintVar(&pad, "pad", defaultPadByte, "Padding byte")
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-pad <pad byte>] file1 file2\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	file1 := loadIHex(flag.Arg(0))
	file2 := loadIHex(flag.Arg(1))

	// Determine lowest and highest addresses in both files to determine
	// comparison range
	start := ^uint32(0)
	end := uint32(0)
	for _, seg := range file1.GetDataSegments() {
		if seg.Address < start {
			start = seg.Address
		}
		if seg.Address + uint32(len(seg.Data)) > end {
			end = seg.Address + uint32(len(seg.Data))
		}
	}
	for _, seg := range file2.GetDataSegments() {
		if seg.Address < start {
			start = seg.Address
		}
		if seg.Address + uint32(len(seg.Data)) > end {
			end = seg.Address + uint32(len(seg.Data))
		}
	}
	fmt.Printf("Comparing 0x%08x to 0x%08x\n", start, end)

	// Compare decoded binary data from both files
	eq := bytes.Equal(file1.ToBinary(start, end, byte(pad)), file2.ToBinary(start, end, byte(pad)))
	if !eq {
		fmt.Println("Files differ")
		os.Exit(1)
	}

	fmt.Println("Files are equal")
	os.Exit(0)
}

func loadIHex(filename string) *gohex.Memory {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	mem := gohex.NewMemory()
	err = mem.ParseIntelHex(file)
	if err != nil {
		panic(err)
	}

	return mem
}
