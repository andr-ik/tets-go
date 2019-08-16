run: _run
build: _build
parse: _parse
open: _open

_run:
	docker-compose -f ~/Propeller/my/tets-php/docker-compose.yml up -d
	go run server/main.go

_build:
	cd ~/pn && ng serve --open

_parse:
	go run main.go --reload
	open http://127.0.0.1:4200

_open:
	open http://127.0.0.1:4200