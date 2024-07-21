
run: clean build
	BLOG_WEB_ENV=dev BLOG_WEB_PORT=:8080 air

deploy: build
	gcloud app deploy --project buildmage --verbosity=info

clean:
	rm -Rf output

build:
	rm -Rf output/static
