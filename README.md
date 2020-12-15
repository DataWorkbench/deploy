# Deploy project for DataWorkbench


Develop all service of the data workbench at local that base on docker-compose;


## Images

#### datawh/builder

The base image used to build other images or compile service;

#### datawh/datawh

The main image to run service with docker-compose;

#### datawh/flyway

The image used to migrate database; If you update the table struct or data, 
run `make compose-migrate-db` to make it works at local develop environment;


## Develop Process

Running services of DataWorkbench at local, you need:

pull all services code of DataWorkbench under same directory；

install docker-compose at local；

all command in this section execute under project `deploy`;


- pull builder

make update-builder

- build all images

make build-all

- Launch dataworkbench services at local

make compose-up

- check logs of service

make compose-logs-f [service=apiserver,spacemanager]


After all services running, you could write code, then:

run `make compose-migrate-db` to update the database if needed;

run `make update [service=apiserver]` to update the service;
