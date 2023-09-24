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
./simpel-api -m=migrate
```

This command for run api :

```shell
./simpel-api 
```

### RUN WITH DOCKER

This command build and run container `api` simpel service

```shell
// build process
$ docker build --rm --tag simpel-api:latest -f Dockerfile .
// run
$ docker run --rm -p 9016:9016 --name simpel-api simpel-api:latest
```
```shell
// run docker with local port
$ docker run --rm --net=host -d -p 8080:8080 --name simpel-api simpel-api:latest
```