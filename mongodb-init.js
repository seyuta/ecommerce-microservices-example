db.createUser(
  {
      user: "example",
      pwd: "123",
      roles: [
          {
              role: "readWrite",
              db: "ecommerce-microservices-example"
          }
      ]
  }
);