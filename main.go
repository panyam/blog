package main

import (
	"bytes"
	"errors"
	"flag"
	"log"
	"strings"
	"text/template"

	"github.com/panyam/blog/web"
)

var (
	gw_addr = flag.String("gw_addr", web.DefaultGatewayAddress(), "Address where the http grpc gateway endpoint is running")
)

func main() {
	flag.Parse()

	// playaround()
	if true {
		ohweb, err := web.NewWebApp(*gw_addr)
		if err != nil {
			log.Fatal(err)
		}
		ohweb.Start()
	}
}

func playaround() {
	t := template.New("")
	t.Funcs(map[string]any{
		"section": func(name string) any {
			return 3
		},
		"use": func(name string) any {
			return 3
		},

		"ucase": func(val string) string { return strings.ToUpper(val) },
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			log.Println("Values For Dict: ", values)
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	})
	t.Parse(`
		{{ define "BaseContents" }}
			{{ .Children }}
			Params - {{ .Params }}
		{{ end }}
		{{ define "Base" }}
			{{ $Children := ucase "hello world" }}
			{{ $Year := "2024" }}
			{{ $ChildParams := dict "A" $Year }}
			{{ template "BaseContents" dict "Children" $Children "Params" $ChildParams "X" (dict "a" 2) }}
		{{ end }}
		{{ template "Base" }}
		`)
	b := bytes.NewBufferString("")
	err := t.Execute(b, map[string]any{})
	log.Println("Err: ", err)
	log.Println("Output: ", b)

	/*
		ts := make(map[string]*tparse.Tree)
		fns := map[string]any{
			"hello": func(name string) any {
				return 3
			},
			"use": func(name string) any {
				return 3
			},
		}
		t1 := &tparse.Tree{}
		t2, err := t1.Parse(`{{ hello }} world {{ end }}`, "{{", "}}", ts, fns)
		log.Println("T1: ", t1)
		log.Println("T2: ", t2, err)
	*/
}
