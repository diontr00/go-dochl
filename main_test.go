package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	. "github.com/diontr00/go-dochl/internal"
)

func captureOutput(f func()) string {
	var buf = new(bytes.Buffer)
	log.SetFlags(0)
	log.SetOutput(buf)
	f()
	log.SetOutput(os.Stdout)
	return buf.String()
}

func TestParser(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		keyswords []string
		expect    []string
	}{
		{
			name:      "file with default keywords",
			path:      "test/test1_test.go",
			keyswords: nil,
			expect: []string{
				`test/test1_test.go:5:TODO(fix): need to change x to 2`,
				`test/test1_test.go:8:HACK: modified to hack around`,
			},
		},

		{
			name:      "directory with default keywords ",
			path:      "./test",
			keyswords: nil,
			expect: []string{
				`test/test1_test.go:5:TODO(fix): need to change x to 2`,
				`test/test1_test.go:8:HACK: modified to hack around`,
				`test/nested/test2_test.go:5:FIXME: This will break`,
			},
		},

		{
			name:      "file with custom keywords ",
			path:      "test2/test3_test.go",
			keyswords: []string{"warning", "hey"},
			expect: []string{
				`test2/test3_test.go:5:HEY(fix): look how cool`,
				`test2/test3_test.go:8:WARNING: have been modified`,
			},
		},

		{
			name:      "directory with custom keywords ",
			path:      "test2/",
			keyswords: []string{"warning", "hey", "error"},
			expect: []string{
				`test2/test3_test.go:5:HEY(fix): look how cool`,
				`test2/test3_test.go:8:WARNING: have been modified`,
				`test2/nested/sub/nested/test4_test.go:5:ERROR: Will be error`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			d := &dochl{sync.WaitGroup{}}
			if tt.keyswords != nil {
				SetKeyWords(tt.keyswords)
			}
			got := captureOutput(func() {
				path, err := os.Stat(tt.path)
				if err != nil {
					t.Errorf("Cannot process path %s : %v", tt.path, err)
				}

				if path.IsDir() {
					d.wg.Add(1)
					d.parseDir(tt.path)
					d.wg.Done()
					d.wg.Wait()
				} else {
					d.parseFile(tt.path)
				}
			})

			var b bytes.Buffer
			for _, e := range tt.expect {
				b.WriteString(e)
				b.WriteString("\n")
			}

			expected := b.String()
			gots := strings.Split(got, "\n")

			if len(tt.expect) != len(gots)-1 {
				t.Fatalf("Expect %d result , but got %d  result", len(tt.expect), len(gots)-1)
			}

			for _, g := range gots {
				if !strings.Contains(expected, g) {
					t.Fatalf("Missing %s", got)
				}
			}
		})
	}
}
