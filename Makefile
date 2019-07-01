fmt:
		go fmt github.com/MonsantoCo/mocka/...

test:
		go test github.com/MonsantoCo/mocka/... -v -cover -coverprofile=coverage.out

coverage: test
		go tool cover -html=coverage.out
