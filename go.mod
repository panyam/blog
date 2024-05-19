module github.com/panyam/blog

go 1.22

// replace github.com/panyam/goutils v0.0.97 => ../../golang/goutils/

require (
	github.com/alexedwards/scs/v2 v2.8.0
	github.com/felixge/httpsnoop v1.0.4
	github.com/gorilla/mux v1.8.1
	github.com/morrisxyang/xreflect v0.0.0-20231001053442-6df0df9858ba
	github.com/panyam/goutils v0.0.97
	github.com/panyam/s3gen v0.0.1
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/adrg/frontmatter v0.2.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/radovskyb/watcher v1.0.7 // indirect
	github.com/yuin/goldmark v1.7.1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace github.com/panyam/goutils v0.0.97 => ../goutils/

replace github.com/panyam/s3gen v0.0.1 => ../s3gen/
