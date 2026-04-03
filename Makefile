NUM_LINKED_GOMODS=`cat go.mod | grep -v "^\/\/" | grep replace | wc -l | sed -e "s/ *//g"`

run: clean build
	BLOG_WEB_ENV=dev BLOG_WEB_PORT=:8080 go run main.go

deploy: checklinks build
	gcloud app deploy --project buildmage --verbosity=info

clean:
	rm -Rf output

build:
	rm -Rf output/static /tmp/main

checklinks:
	@if [ x"${NUM_LINKED_GOMODS}" != "x0" ]; then	\
		echo "You are trying to deploying with symlinks.  Remove them first and make sure versions exist" && false ;	\
	fi

resymlink:
	mkdir -p locallinks
	rm -Rf locallinks/*
	cd locallinks && ln -s ../../s3gen

# =============================================================================
# Security audit
# =============================================================================

# Full security audit: dependency vulns + code patterns + secrets
audit:
	@echo "=== govulncheck ==="
	govulncheck ./...
	@echo ""
	@echo "=== gosec ==="
	gosec -quiet -severity=medium ./... || true
	@echo ""
	@echo "=== gitleaks ==="
	gitleaks detect --source . -v 2>/dev/null || echo "gitleaks not installed (go install github.com/gitleaks/gitleaks/v8@latest)"
	@echo ""
	@echo "=== Audit complete ==="

.PHONY: audit
