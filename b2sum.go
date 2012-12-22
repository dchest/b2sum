// b2sum command calculates BLAKE2 checksums of files.
package main

import (
	"flag"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/dchest/blake2b"
)

var (
	algoFlag = flag.String("a", "blake2b", "hash algorithm")
	sizeFlag = flag.Int("s", 64, "digest size in bytes")
)

func calcSum(f *os.File, h hash.Hash) (sum []byte, err error) {
	h.Reset()
	_, err = io.Copy(h, f)
	sum = h.Sum(nil)
	return
}

func main() {
	flag.Parse()

	if *algoFlag != "blake2b" {
		flag.Usage()
		fmt.Fprintf(os.Stderr, `only "blake2b" algorithm is currently supported`)
		os.Exit(1)
	}

	if *sizeFlag > blake2b.Size {
		fmt.Fprintf(os.Stderr, "error: size too large")
		os.Exit(1)
	}
	h, err := blake2b.New(&blake2b.Config{Size: uint8(*sizeFlag)})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %s", err)
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		// Read from stdin.
		sum, err := calcSum(os.Stdin, h)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%x\n", sum)
		os.Exit(0)
	}
	exitNo := 0
	for i := 0; i < flag.NArg(); i++ {
		filename := flag.Arg(i)
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "(%s) %s\n", filename, err)
			exitNo = 1
			continue
		}
		sum, err := calcSum(f, h)
		f.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "(%s) %s\n", filename, err)
			exitNo = 1
			continue
		}
		fmt.Printf("BLAKE2b-%d (%s) = %x\n", h.Size(), filename, sum)
	}
	os.Exit(exitNo)
}
