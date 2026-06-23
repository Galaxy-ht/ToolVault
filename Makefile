.PHONY: bootstrap-check m1-registry-check test verify

bootstrap-check:
	./scripts/bootstrap-check.sh

m1-registry-check:
	./scripts/m1-registry-check.sh

test:
	go test ./...

verify: bootstrap-check m1-registry-check
