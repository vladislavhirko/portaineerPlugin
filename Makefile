install:
	mkdir -p ~/.portaineerPlugin
	cp example_config/config.toml ~/.portaineerPlugin/config.toml
	export GO111MODULE=on
	go install
