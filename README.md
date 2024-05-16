
# Project Title

A brief description of what this project does and who it's for

### Command for running this program:

Copy file .env.example to .env and edit:
```shell
$ cp .env.example .env
```

After edit file config, build go project:

```shell
go build -v .
```

This command for run migrate database:

```shell
./owlharbour-api -m="migrate"
```
This command for run seed database:

```shell
./owlharbour-api -s="seeder"
```

Default seed credentials :

superadmin@gmail.com :
superadminpassword

admin@gmail.com :
adminpassword

This command for run api :

```shell
./owlharbour-api 
```

### RUN WITH DOCKER

This command build and run container `api` simpel service

```shell
// build process
$ docker build --rm --tag owlharbour-api:latest -f Dockerfile .
// run
$ docker run --rm -p 9016:9016 --name owlharbour-api owlharbour-api:latest
```
```shell
// run docker with local port
$ docker run --rm --net=host -d -p 8080:8080 --name owlharbour-api owlharbour-api:latest
```