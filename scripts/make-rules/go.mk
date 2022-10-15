.PHONY: go.build.%
go.build.%: clean
	@mkdir -p $(OUTPUT)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o $(OUTPUT)/$*.out ./cmd/$*/*.go
