fmt:
		go fmt ./...

vet:
		go vet ./...

lint: fmt
		revive -formatter=stylish -config=revive.toml ./...

test:
		go test ./... -v -cover -coverprofile=coverage.out

test-debug:
		dlv test ./... --log

coverage: test
		go tool cover -html=coverage.out
