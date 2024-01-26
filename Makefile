GONUBOARD_BINARY := gonuboard

.PHONY: build
build:
	@go build -o ${GONUBOARD_BINARY} ./cmd/gonuboard

.PHONY: clean
clean:
	@rm -f ${GONUBOARD_BINARY}
