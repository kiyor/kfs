VERSION   := $(shell cat ./version)

release:
	git tag -a $(VERSION) -m "release" || true
	git push origin main --tags
.PHONY: release
