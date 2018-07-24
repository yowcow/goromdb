all: _watcher _storage main

_%:
	mkdir -p $@

main:
	go build -o $@

clean:
	rm -rf main

realclean: clean
	rm -rf _watcher _storage
