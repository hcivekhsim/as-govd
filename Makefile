.PHONY: serve 

serve:
	CompileDaemon -directory=. -exclude-dir=.database

run:
	docker compose up --build



