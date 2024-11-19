weatherdb: cmd/main.go cmd/measurement.go
	go build -v -o weatherdb cmd/main.go cmd/measurement.go

.PHONY: clean
clean:
	rm weatherdb

.PHONY: run
run: weatherdb
	./weatherdb
