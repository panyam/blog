
run: build
	BLOG_WEB_ENV=dev BLOG_WEB_PORT=:8095 air

deploy: build
	gcloud app deploy --project maniacalbuilder --verbosity=info

build:
	cd web/lib ; npm run build
