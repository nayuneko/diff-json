package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kylelemons/godebug/pretty"
)

func main() {
	diffOnly := flag.Bool("d", false, "diff only")
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "[-d]", "jsonfile1", "jsonfile2")
		fmt.Fprintln(os.Stderr, "  -d diff only")
		os.Exit(1)
	}
	filename1 := flag.Args()[0]
	f1, err := os.Open(filename1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't Open file", filename1)
		os.Exit(1)
	}
	filename2 := flag.Args()[1]
	f2, err := os.Open(filename2)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't Open file", filename2)
		os.Exit(1)
	}
	c, err := diff(f1, f2)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if c != "" {
		modTime := func(f *os.File) string {
			s, err := f.Stat()
			if err != nil {
				return ""
			}
			return s.ModTime().String()
		}
		fmt.Printf("\x1b[31m--- %s %s\n\x1b[0m", f1.Name(), modTime(f1))
		fmt.Printf("\x1b[34m+++ %s %s\n\x1b[0m", f2.Name(), modTime(f2))
		printDiff(c, *diffOnly)
		os.Exit(1)
	}
	fmt.Println("no diff")
}

func diff(file1, file2 *os.File) (string, error) {
	var json1 interface{}
	if err := json.NewDecoder(file1).Decode(&json1); err != nil {
		return "", fmt.Errorf("Failed json decode: %s", file1.Name())
	}
	var json2 interface{}
	if err := json.NewDecoder(file2).Decode(&json2); err != nil {
		return "", fmt.Errorf("Failed json decode: %s", file2.Name())
	}
	return pretty.Compare(json1, json2), nil
}

func printDiff(c string, diffOnly bool) {
	for _, s := range strings.Split(c, "\n") {
		cc := s[:1]
		if cc == "-" {
			fmt.Printf("\x1b[31m%s\n\x1b[0m", s)
		} else if cc == "+" {
			fmt.Printf("\x1b[34m%s\n\x1b[0m", s)
		} else if !diffOnly {
			fmt.Println(s)
		}
	}
}
