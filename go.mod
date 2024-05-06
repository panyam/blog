module github.com/panyam/blog

go 1.22

// replace github.com/panyam/goutils v0.0.97 => ../../golang/goutils/

require (
	github.com/alexedwards/scs/v2 v2.8.0
	github.com/felixge/httpsnoop v1.0.4
	github.com/gorilla/mux v1.8.1
	github.com/panyam/s3gen v0.0.1
)

require github.com/panyam/goutils v0.0.97 // indirect

replace github.com/panyam/goutils v0.0.97 => ../goutils/

replace github.com/panyam/s3gen v0.0.1 => ../s3gen/
