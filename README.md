# Microservices Docker
[Microservices Docker](https://github.com/ariefnhidayah/microservices-docker/) is a repository for simple project microservices using docker. In this repository, I using example 2 services, book service and order service. Those services is a simple service. And I use an api gateway for connect from frontend to services.

Those services using [GO Language](https://golang.org) programming and for api gateway I use [Express JS](https://expressjs.com/). And I use [PostgreSQL](https://www.postgresql.org/) for database system. This repository made for my notes and others who want to learn microservices with docker.

## Pre-Installation
Before you running this project, you must clone this repository to your local. And you must install docker in your pc, because in this repo use docker for containerize those services and the api gateway.

For clone this repo, you can use this command.
```
git clone https://github.com/ariefnhidayah/microservices-docker.git
```

You can install docker use this link \
[Install Docker](https://www.docker.com/get-started)

## Installation
If you already clone this repo and install docker, you can use this command for running this project in your local-docker.
```
docker-compose up -d
```
