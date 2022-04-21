# Auth Service
protogen-auth:
	@rm -Rf s-auth/pkg/pb
	@protoc --proto_path=s-auth/protobuf --go_out=plugins=grpc:. s-auth/protobuf/*.proto
	@echo generating auth service protobuf implementation done...

# Catalog Service
protogen-catalog:
	@rm -Rf s-catalog/pkg/pb
	@protoc --proto_path=s-catalog/protobuf --go_out=plugins=grpc:. s-catalog/protobuf/*.proto
	@echo generating catalog service protobuf implementation done...

# Transactions Service
protogen-transactions:
	@rm -Rf s-transactions/pkg/pb
	@protoc --proto_path=s-transactions/protobuf --go_out=plugins=grpc:. s-transactions/protobuf/*.proto
	@echo generating transactions service protobuf implementation done...

# Deployment
dockerize-auth:
	@docker build --build-arg server=s-auth -f ./s-auth/Dockerfile -t s-auth:latest .
	@docker images --filter=reference='s-auth:latest'

dockerize-catalog:
	@docker build --build-arg server=s-catalog -f ./s-catalog/Dockerfile -t s-catalog:latest .
	@docker images --filter=reference='s-catalog:latest'

dockerize-transactions:
	@docker build --build-arg server=s-transactions -f ./s-transactions/Dockerfile -t s-transactions:latest .
	@docker images --filter=reference='s-transactions:latest'

dockerize-all-service: dockerize-auth dockerize-catalog dockerize-transactions

start-ecommerce-microservices-example-docker:
	@docker-compose up -d

stop-ecommerce-microservices-example-docker:
	@docker-compose down

build-run: dockerize-all-service start-ecommerce-microservices-example-docker