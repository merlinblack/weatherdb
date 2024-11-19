weatherdb: cmd/main.go
	go build -v -o weatherdb cmd/main.go

.PHONY: clean
clean:
	rm weatherdb

.PHONY: run
run: weatherdb
	./weatherdb
