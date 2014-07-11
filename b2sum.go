// b2sum command calculates BLAKE2 checksums of files.
package main

import (
	"flag"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/dchest/blake2b"
	"github.com/dchest/blake2s"
)

var (
	algoFlag = flag.String("a", "blake2b", "hash algorithm (blake2b, blake2s)")
	sizeFlag = flag.Int("s", 0, "digest size in bytes (0 = default)")
)

type hashDesc struct {
	name    string
	maxSize int
	maker   func(size uint8) (hash.Hash, error)
}

var algorithms = map[string]hashDesc{
	"blake2b": {
		"BLAKE2b",
		blake2b.Size,
		func(size uint8) (hash.Hash, error) { return blake2b.New(&blake2b.Config{Size: size}) },
	},
	"blake2s": {
		"BLAKE2s",
		blake2s.Size,
		func(size uint8) (hash.Hash, error) { return blake2s.New(&blake2s.Config{Size: size}) },
	},
}

func calcSum(f *os.File, h hash.Hash) (sum []byte, err error) {
	h.Reset()
	_, err = io.Copy(h, f)
	sum = h.Sum(nil)
	return
}

func main() {
	flag.Parse()

	algo, ok := algorithms[*algoFlag]
	if !ok {
		flag.Usage()
		fmt.Fprintf(os.Stderr, `unsupported algorithm: %s`, *algoFlag)
		os.Exit(1)
	}
	if *sizeFlag == 0 {
		*sizeFlag = algo.maxSize
	} else if *sizeFlag > algo.maxSize {
		fmt.Fprintf(os.Stderr, "error: size too large")
		os.Exit(1)
	}

	h, err := algo.maker(uint8(*sizeFlag))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
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
		fmt.Printf("%s-%d (%s) = %x\n", algo.name, h.Size(), filename, sum)
	}
	os.Exit(exitNo)
}
