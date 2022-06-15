image-tag="latest"

ui:
	cd frontend && npm run build
.PHONY: ui

embed-ui: 
	cd frontend && go-bindata -fs -prefix "./dist/" -pkg "main" -o "../backend/bindata.go" -ignore "\\.go" -ignore "\\.DS_Store" ./dist/...
.PHONY: embed-ui

build-image:
	GORELEASER_CURRENT_TAG=${image-tag} goreleaser build --rm-dist --snapshot
.PHONY: build-image

push-image:
	GORELEASER_CURRENT_TAG=${image-tag} goreleaser release --rm-dist --skip-validate
.PHONY: push-image

release: ui embed-ui build-image push-image
.PHONY: release