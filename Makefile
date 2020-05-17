.PHONY: install
install:
	mkdir -p $(HOME)/.config/octant/plugins/
	rm -f $(HOME)/.config/octant/plugins/octant-velero-plugin
	go build -o $(HOME)/.config/octant/plugins/octant-velero-plugin ./...
