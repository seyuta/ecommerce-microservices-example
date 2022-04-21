## ecommerce-microservices-example
Is a microservice for simple e-commerce, solutions stack with:
1. [Golang](https://go.dev/)
2. [MongoDB](https://www.mongodb.com/)
3. [gRPC](https://grpc.io/)
##

### Prerequisite
1. [BloomRPC ](https://github.com/bloomrpc/bloomrpc) for hit gRPC services.
2. [Docker & Docker Compose](https://docs.docker.com/engine/) installed.
##

### Dockerize
>Build docker image all services:

```bash
docker build --build-arg server=auth -f ./s-auth/Dockerfile -t auth:latest .
docker images --filter=reference='auth:latest'

docker build --build-arg server=catalog -f ./s-catalog/Dockerfile -t catalog:latest .
docker images --filter=reference='catalog:latest'

docker build --build-arg server=transactions -f ./s-transactions/Dockerfile -t transactions:latest .
docker images --filter=reference='transactions:latest'
```

### Run & Stop
>Run all services & database with docker-compose:
```bash
docker-compose up -d
```
>Stop all services & database with docker-compose:
```bash
docker-compose down
```

### Simplify Build & Run
>Before execute, make sure your operating system has installed [Makefile](https://makefiletutorial.com/).
```bash
make build-run
```
##

### System Diagram
![diagram](https://github.com/seyuta/ecommerce-microservices-example/blob/master/diagram.jpg?raw=true)

### Usage BloomRPC
>Open BloomRPC and import proto file from each service: `./{service path}/protobuf`
##### Example:
![bloomrpc-login](https://github.com/seyuta/ecommerce-microservices-example/blob/master/usage-bloomrpc/bloomrpc-login.jpg?raw=true)
##### Example With Auth:
![bloomrpc-order](https://github.com/seyuta/ecommerce-microservices-example/blob/master/usage-bloomrpc/bloomrpc-order.jpg?raw=true)
##### Format Metadata For Auth:
```bash
{"Authorization": "Bearer paste_token_here"}
```
##

### Testing
>Run Go unit test.
```bash
go test -v ./{service path}/service
```

### Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

### License
[MIT](https://choosealicense.com/licenses/mit/)