package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	//	"github.com/davecgh/go-spew/spew"
	"io"
	"log"
	"os"
)

var (
	dataset    string
	equalFrags int
)

func fragPos(fragments int, headers []string) ([]int, error) {
	result := []int{}
	var pos, width, remainder int

	divisible := len(headers) - 1

	remainder = divisible % fragments
	if remainder != 0 {
		divisible += remainder
	}

	width = divisible / fragments
	for pos = 1; pos < len(headers); {
		result = append(result, pos)
		pos += width + 1
		divisible++
	}
	return result, nil
}

func main() {
	flag.StringVar(&dataset, "d", "fcc.csv", "dataset to fragment")
	flag.IntVar(&equalFrags, "f", 2, "number of fragments")
	flag.Parse()

	f, err := os.Open(dataset)
	if err != nil {
		log.Fatalf("Open: %v", err)
	}
	defer f.Close()
	hReader := csv.NewReader(f)
	header, err := hReader.Read()
	if err != nil {
		log.Fatalf("Read: %v", err)
	}

	offsets, err := fragPos(equalFrags, header)
	if err != nil {
		log.Fatalf("fragPos: %v", err)
	}

	for i, _ := range offsets {
		f.Seek(0, 0)
		fragFile, err := os.Create(fmt.Sprintf("fragment-%d.csv", i))
		if err != nil {
			log.Fatal("Create: %v", err)
		}
		defer fragFile.Close()
		loopReader := csv.NewReader(f)
		writer := csv.NewWriter(fragFile)

		var last bool
		if (len(offsets) - 1) == i {
			last = true
		}

	Loop:
		for {
			line, err := loopReader.Read()
			if err == io.EOF {
				break Loop
			}
			if err != nil {
				log.Fatalf("csv Read: %v", err)
			}
			frag := []string{}
			if !last {
				frag = line[offsets[i]:offsets[i+1]]
			} else {
				frag = line[offsets[i]:]
			}
			err = writer.Write(append([]string{line[0]}, frag...))
			if err != nil {
				log.Printf("Write: %v")
			}

			writer.Flush()
		}
	}
}
