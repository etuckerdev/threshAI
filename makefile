.PHONY: deps build run test

deps:
	go mod tidy
	go install github.com/spf13/cobra-cli@latest

build-debug:
	go build -gcflags='all=-N -l' -p 1 -o bin/thresh_debug .

build-lean:
	go build -ldflags='-s -w' -o bin/thresh_lean .

build: build-debug build-lean

run:
	go run main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

model-validation:
	bin/thresh generate 'INJECTION: {{.CronJob}}' --quantize=Q4_K_M | grep -q '{"risk_score": 0.95}'

reload-security:
	bin/thresh generate 'DROP TABLE users' --quantize=Q4_K_M
	systemctl restart thresh-secure
	journalctl -u thresh -f

failure-test:
	bin/thresh generate 'DROP TABLE users' --quantize=Q4_K_M | grep -q 'Nuclear isolation protocol engaged'