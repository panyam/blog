
run: build
	BLOG_WEB_ENV=dev BLOG_WEB_PORT=:8080 air

deploy: build
	gcloud app deploy --project buildmage --verbosity=info

build:
	rm -Rf output/static
