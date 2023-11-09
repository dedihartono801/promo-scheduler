GOPATH ?= $(HOME)/go

run-scheduler:
	go run cmd/scheduler/main.go

run-worker:
	go run cmd/worker/main.go

format:
	go fmt ./...

scheduler:
	DC_APP_ENV=dev $(GOPATH)/bin/reflex -s -r '\.go$$' make format run-scheduler

worker:
	DC_APP_ENV=dev $(GOPATH)/bin/reflex -s -r '\.go$$' make format run-worker