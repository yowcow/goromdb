ARTIFACTS = stress-bolt stress-bolt-ns

all: _watcher _storage $(ARTIFACTS)

_%:
	mkdir -p $@

stress-%:
	go build -o $@ ./$*

clean:
	rm -rf $(ARTIFACTS)

realclean: clean
	rm -rf _watcher _storage

.PHONY: all clean realclean
