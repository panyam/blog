NUM_LINKED_GOMODS=`cat go.mod | grep -v "^\/\/" | grep replace | wc -l | sed -e "s/ *//g"`

run: clean build
	BLOG_WEB_ENV=dev BLOG_WEB_PORT=:8080 air

deploy: checklinks build
	gcloud app deploy --project buildmage --verbosity=info

clean:
	rm -Rf output

build:
	rm -Rf output/static

checklinks:
	@if [ x"${NUM_LINKED_GOMODS}" != "x0" ]; then	\
		echo "You are trying to deploying with symlinks.  Remove them first and make sure versions exist" && false ;	\
	fi

resymlink:
	mkdir -p locallinks
	rm -Rf locallinks/*
	cd locallinks && ln -s ../../s3gen

files4prompt:
	@find . | grep -v "\.git" | grep -v "\.sh" | grep -v "\.css" | grep -v "\.png" | grep -v "\.jpg" | grep -v drawio | grep -v "\.svg" | grep -v "\.jpeg" | grep -v output | grep -v "\.env" | grep -v vscode

prompt:
	source ~/personal/.shhelpers && files_for_llm `make files4prompt`
