install:
	mkdir -p ~/.portaineer_plugin
	cp example_config/config.toml ~/.portaineer_plugin/config.toml
	export GO111MODULE=on
	go install
