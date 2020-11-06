MIN_COVERAGE = 70
test:
	go test ./... -race -count=1 -cover -coverprofile=coverage.txt && go tool cover -func=coverage.txt \
	| grep total | sed 's/\%//g' | awk '{err=0}{if ($$3 > 0 && $$3 < $(MIN_COVERAGE)) {printf "=== FAIL: Coverage failed at %.2f%%\n", $$3; err=1}} END {exit err}'

lint:
	golangci-lint run --deadline=5m -v

build:
	go build ./cmd/api/

install_moq:
	go get github.com/matryer/moq

up:
	docker-compose -f docker/docker-compose.yml up -d --build

down:
	docker-compose -f docker/docker-compose.yml down
