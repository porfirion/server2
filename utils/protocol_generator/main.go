package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"text/template"

	"github.com/porfirion/server2/messages"
)

var (
	flagOutputFilename string

	tpl = template.Must(template.New("").Parse(``))
)

func init() {
	flag.StringVar(&flagOutputFilename, "out", "", "")
	flag.Parse()
}

func main() {
	if flagOutputFilename == "" {
		log.Fatal("output file not specified")
	}

	out, err := os.Create(flagOutputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	defer w.Flush()

	for tp, id := range messages.AvailableMessageTypes {
		log.Printf("%d: %+v\n", id, tp)
		tpl.Execute(w, tp)
	}
}
