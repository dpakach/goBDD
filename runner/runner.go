package runner


import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dpakach/goBDD/suite"


	"github.com/dpakach/gorkin/lexer"
	"github.com/dpakach/gorkin/parser"
)

func Run(s *suite.Suite) {
	if len(os.Args) < 2 {
		log.Fatal(fmt.Errorf("Opps, Seems like you forgot to provide the path of the feature file"))
		os.Exit(1)
	}
	path := os.Args[1]
	abs, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(fmt.Errorf("Invalid path provided, Make sure the path %q is correct", path))
		os.Exit(1)
	}
	fi, err := os.Stat(abs)
	if os.IsNotExist(err) {
		log.Fatal(fmt.Errorf("Error, Make sure the path %q exists", path))
		os.Exit(1)
	}

	out := new(bytes.Buffer)
	start := time.Now()
	switch mode := fi.Mode(); {
	case mode.IsDir():
		err := filepath.Walk(abs, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if ext := filepath.Ext(info.Name()); ext == ".feature" {
					l := lexer.NewFromFile(path)
					p := parser.New(l)
					ps := p.Parse()
					s.Run(ps)
				}
			}
			return nil
		})
		if err != nil {
			log.Println(err)
		}
	case mode.IsRegular():
		l := lexer.NewFromFile(path)
		p := parser.New(l)
		ps := p.Parse()
		s.Run(ps)
	}
	io.Copy(os.Stdout, out)

	io.WriteString(os.Stdout, fmt.Sprintf("\nCompleted in %v\n", time.Since(start)))
}
