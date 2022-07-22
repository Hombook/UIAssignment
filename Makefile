.DEFAULT_GOAL := help

postgres_container := postgresql

#help:	@ List available tasks on this project
help:
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#'  | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
#start_db: @ Startup PostgreSQL DB
start_db:
	docker pull postgres
	docker run --name $(postgres_container) -e POSTGRES_USER=ui_test -e POSTGRES_DB=ui_test -e POSTGRES_PASSWORD=uiPassword5678 -p 5432:5432 -d postgres
#init_db: @ Populate DB with users table
init_db:
	timeout 90s bash -c "until docker exec $(postgres_container) pg_isready ; do sleep 5 ; done"
	timeout 90s bash -c "until docker exec $(postgres_container) psql -h 127.0.0.1 -U ui_test -d ui_test ; do sleep 5 ; done"
	docker cp ./ui_test.sql $(postgres_container):/ui_test.sql
	docker exec $(postgres_container) psql -h 127.0.0.1 -U ui_test -d ui_test -f /ui_test.sql
#clean: @ Stop container and remove generated files
clean:
	- docker stop $(postgres_container)
	- docker rm $(postgres_container)
#run: @ Run the full startup scripts
run: start_db init_db
