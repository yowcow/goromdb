all: _watcher _storage main

_%:
	mkdir -p $@

main:
	go build -o $@

clean:
	rm -rf main
