package main

import (
	"flag"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	. "github.com/diontr00/go-dochl/internal"
)

var (
	// keys word to look up
	fset   = token.NewFileSet()
	kwflag = flag.String("keys", "todo,fixme,bug,hack", "keywords to be extracted")
)

type dochl struct {
	wg sync.WaitGroup
}

func (d *dochl) parseFile(path string) error {
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		NewError(path, err)
	}
	for _, c := range f.Comments {
		for _, ci := range c.List {
			NewComment(fset, ci).Parse()
		}
	}
	return nil
}

func (d *dochl) parseDir(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		NewError(path, err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			subpath := filepath.Join(path, file.Name())

			d.wg.Add(1)
			go func(subpath string) {
				defer d.wg.Done()
				d.parseDir(subpath)
			}(subpath)
		}
	}

	f, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		NewError(path, err)
		return
	}

	for _, pkg := range f {
		for _, file := range pkg.Files {
			for _, c := range file.Comments {
				for _, ci := range c.List {
					NewComment(fset, ci).Parse()
				}
			}
		}
	}
	return
}

func main() {
	flag.Parse()
	keywords := strings.Split(*kwflag, ",")
	if len(keywords) > 0 {
		SetKeyWords(keywords)
	}

	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	d := &dochl{wg: sync.WaitGroup{}}

	if len(flag.Args()) == 0 {
		d.parseDir(".")
		os.Exit(0)
	}

	for _, path := range flag.Args() {
		d.wg.Add(1)
		go func(path string) {
			defer d.wg.Done()

			stat, err := os.Stat(path)
			if err != nil {
				NewError(path, err)
				return
			}

			if stat.IsDir() {
				d.parseDir(path)
			} else {
				d.parseFile(path)
			}
		}(path)
	}

	d.wg.Wait()
}
