.DEFAULT_GOAL := help

docker_network := uiassignment-network
postgres_container := postgresql
uiassignment_container := uiassignment

#help:	@ List available tasks on this project
help:
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#'  | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
#setup_docker_network: @ Setup Docker network for inter-container connection
setup_docker_network:
	docker network create $(docker_network)
#start_db: @ Startup PostgreSQL DB
start_db:
	docker pull postgres
	docker run --name $(postgres_container) -e POSTGRES_USER=ui_test -e POSTGRES_DB=ui_test -e POSTGRES_PASSWORD=uiPassword5678 -p 5432:5432 -d postgres
	docker network connect $(docker_network) $(postgres_container)
#init_db: @ Populate DB with users table
init_db:
	timeout 90s bash -c "until docker exec $(postgres_container) pg_isready ; do sleep 5 ; done"
	timeout 90s bash -c "until docker exec $(postgres_container) psql -h 127.0.0.1 -U ui_test -d ui_test ; do sleep 5 ; done"
	docker cp ./db/ui_test.sql $(postgres_container):/ui_test.sql
	docker exec $(postgres_container) psql -h 127.0.0.1 -U ui_test -d ui_test -f /ui_test.sql
#build: @ Build UI assignment REST service Docker image
build:
	docker build -t uiassignment .
#start_server: @ Start UI assignment REST service
start_server:
	docker run --name $(uiassignment_container) -p 80:80 -d uiassignment
	docker network connect $(docker_network) $(uiassignment_container)
#stop_server: @ Stop UI assignment REST service
stop_server:
	- docker network disconnect $(docker_network) $(uiassignment_container)
	- docker stop $(uiassignment_container)
	- docker rm $(uiassignment_container)
#clean: @ Stop container, network and remove built images
clean: stop_server
	- docker network disconnect $(docker_network) $(postgres_container)
	- docker stop $(postgres_container)
	- docker rm $(postgres_container)
	- docker network rm $(docker_network)
	- docker rmi uiassignment
	- docker image prune -f
#run: @ Run the full startup scripts
run: setup_docker_network start_db init_db start_server
