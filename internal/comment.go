package internal

import (
	"bufio"
	"bytes"
	"go/ast"
	"go/token"
	"log"
	"sort"
	"strings"
)

var (
	keywords = keys{"todo", "fixme", "bug", "hack"}
)

func SetKeyWords(kw []string) {
	keywords = keys(kw)
}

type keys []string

func (k keys) sort() {
	sort.Slice(k, func(i, j int) bool {
		return len(k[i]) < len(k[j])
	})
}

// get the minimum len of keys to prevent redundant look up
func (k keys) minL() int {
	k.sort()
	return len(k[0])
}

// Represent comment of underlying *comment node
type comment struct {
	fs       *token.FileSet
	comment  *ast.Comment
	lineBuf  *bufio.Reader
	loffset  int
	tokenPos token.Position
}

// check *comment position, look up fileset if invalid
func (c *comment) pos() *token.Position {
	if !c.tokenPos.IsValid() {
		c.tokenPos = c.fs.Position(c.comment.Pos())
	}
	return &c.tokenPos
}

// get the file path of *comment
func (c *comment) path() string {
	return c.pos().Filename
}

// highlight print *comment if contain keywords
func (c *comment) Parse() {
	var offset = 0
	for {
		line, _, err := c.lineBuf.ReadLine()
		if err != nil {
			break
		}
		cmt := string(bytes.TrimSpace(line))
		if len(cmt) < keywords.minL() {
			offset++
			continue
		}
		for _, kw := range keywords {
			if strings.EqualFold(kw, string(cmt[0:len(kw)])) {
				log.Printf("%s:%s:%s%s\n", Green(c.path()), Yellow(c.pos().Line+c.loffset), Red(strings.ToUpper(kw), B), cmt[len(kw):])
			}
		}
		c.loffset++
	}
	return
}

// generate new  *comment base on the  Comment node
func NewComment(fs *token.FileSet, c *ast.Comment) *comment {
	t := c.Text
	cmt := strings.TrimSpace(t)
	switch string(cmt[1]) {
	case "/":
		cmt = cmt[2:]
	case "*":
		cmt = cmt[2 : len(cmt)-2]
	}

	return &comment{
		fs:      fs,
		comment: c,
		lineBuf: bufio.NewReader(bytes.NewBufferString(cmt)),
	}
}
