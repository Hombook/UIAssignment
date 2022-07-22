.DEFAULT_GOAL := help

postgres_container := postgresql

#help:	@ List available tasks on this project
help:
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#'  | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
#start_db: @ Startup PostgreSQL DB
start_db:
	docker pull postgres
	docker run --name $(postgres_container) -e POSTGRES_USER=ui_test -e POSTGRES_DB=ui_test -e POSTGRES_PASSWORD=uiPassword5678 -p 5432:5432 -d postgres
#clean: @ Stop container and remove generated files
clean:
	- docker stop $(postgres_container)
	- docker rm $(postgres_container)
#run: @ Run the full startup scripts
run: start_db 
