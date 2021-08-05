package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Node interface {
	fmt.Stringer
}

type Directory struct {
	name     string
	children []Node
}

type File struct {
	name string
	size int64
}

func (f File) String() string {
	if f.size == 0 {
		return f.name + " (empty)"
	}

	return f.name + " (" + strconv.Itoa(int(f.size)) + "b)"
}

func (d Directory) String() string {
	return d.name
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err, nodes := readDir(path, nil, printFiles)
	printDir(out, nodes, nil)

	return err
}

func readDir(path string, nodes []Node, printFiles bool) (error, []Node) {
	dirsAndFiles, err := ioutil.ReadDir(path)

	sort.Slice(dirsAndFiles, func(i, j int) bool {
		return dirsAndFiles[i].Name() < dirsAndFiles[j].Name()
	})

	for _, info := range dirsAndFiles {
		if info.IsDir() || printFiles {
			var newNode Node

			if info.IsDir() {
				_, children := readDir(filepath.Join(path, info.Name()), nil, printFiles)
				newNode = Directory{info.Name(), children}
			} else {
				newNode = File{info.Name(), info.Size()}
			}

			nodes = append(nodes, newNode)
		}
	}

	return err, nodes
}

func printDir(out io.Writer, nodes []Node, prefix []string) {
	if len(nodes) == 0 {
		return
	}

	fmt.Fprintf(out, "%s", strings.Join(prefix, ""))

	node := nodes[0]

	if len(nodes) == 1 {
		fmt.Fprintf(out, "%s%s\n", "└───", node)

		if directory, ok := node.(Directory); ok {
			printDir(out, directory.children, append(prefix, "\t"))
		}

		return
	}

	fmt.Fprintf(out, "%s%s\n", "├───", node)
	if directory, ok := node.(Directory); ok {
		printDir(out, directory.children, append(prefix, "│\t"))
	}

	printDir(out, nodes[1:], prefix)
}
