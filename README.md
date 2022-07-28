# Overview
The general purpose of this project is to fulfil the tasks from UI's assignment.

# What It Does
1. Start a non-persistent PostgreSQL container with DB: ui_test
2. Populate the DB with users table
3. Start a REST service container with these APIs
    - GET /v1/users
    - GET /v1/users/<account>
    - POST /v1/users
    - POST /v1/accessToken
    - DELETE /v1/users/<account>
    - PATCH /v1/users/<account>

# How To Use
## Prerequisite
* A Linux OS with Docker engine and docker-compose installed.
  - Verified in CentOS 7.8
  - Refer to the official installation guide: https://docs.docker.com/engine/install/

## Quick Start
### The makefile way
To start up the services
<pre><code>make build
make run</code></pre>
To shutdown everything and cleanup built images
<pre><code>make clean</code></pre>
#### Optional
To stop the API server only
<pre><code>make stop_server
</code></pre>
To start the API server
<pre><code>make start_server
</code></pre>
Environment variables that can be used for the API server:
<pre><code>(env variable)=(default value)
POSTGRES_HOST=postgresql
POSTGRES_PORT=5432
POSTGRES_USER=ui_test
POSTGRES_PWD=iPassword5678
</code></pre>
### The docker-compose way
To start up the services:
<pre><code>docker compose up
</code></pre>
or run it in background
<pre><code>docker compose up -d
</code></pre>
To shut down the services
<pre><code>docker compose stop
</code></pre>
To cleanup the containers
<pre><code>docker compose rm</code></pre>

## Generate Swagger Doc
Install Swaggo
<pre><code>go install github.com/swaggo/swag/cmd/swag@v1.8.1</code></pre>
Generate Swagger doc
<pre><code>swag init -d ./cmd/uiassignment -o ./docs --parseDependency</code></pre>
* Swagger document can be found under {project root}/docs
* To view the document, paste the content of swagger.yaml to https://editor.swagger.io/

# WebSocket Demo
The web chat interface for the demo can be found under http://{your.IP}/web/chat
* Notification message will be sent when an existing account failed on POST /api/v1/accessToken
  - Only if it fails on password verification.
* The demo is a slightly modified version of https://github.com/gorilla/websocket/tree/master/examples/chat
