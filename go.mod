module github.com/panyam/blog

go 1.22

// replace github.com/panyam/goutils v0.0.97 => ../../golang/goutils/

require (
	github.com/felixge/httpsnoop v1.0.4
	github.com/gorilla/mux v1.8.1
	github.com/panyam/goutils v0.1.1
	github.com/panyam/s3gen v0.0.9
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/adrg/frontmatter v0.2.0 // indirect
	github.com/alecthomas/chroma/v2 v2.14.0 // indirect
	github.com/dlclark/regexp2 v1.11.0 // indirect
	github.com/morrisxyang/xreflect v0.0.0-20231001053442-6df0df9858ba // indirect
	github.com/radovskyb/watcher v1.0.7 // indirect
	github.com/yuin/goldmark v1.7.1 // indirect
	github.com/yuin/goldmark-highlighting/v2 v2.0.0-20230729083705-37449abec8cc // indirect
	go.abhg.dev/goldmark/anchor v0.1.1 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

// replace github.com/panyam/goutils v0.1.1 => ../goutils/
// replace github.com/panyam/s3gen v0.0.9 => ../s3gen/
