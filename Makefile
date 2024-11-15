weatherdb: main.go
	go build -v -o weatherdb main.go

.PHONY: clean
clean:
	rm weatherdb

.PHONY: run
run: weatherdb
	./weatherdb
