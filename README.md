# CRS RabbitMQ Connector

trivial example of populating updates from RabbitMQ message queue into Google Cloud Retail API Product Catalog.

## Usage

1. obtain a service account that has Retail Editor permissions, Manage Keys > Create a Key, download the key and move it to root of the repository named as `service-account.json`
2. run the RabbitMQ instance

```sh
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management
```

3. run the subscriber

```sh
go run . --subscriber
```

from here on the body of any messages from the queue named `hello` will be used to set the title of the product with ID `test-product`

4. run the publisher

```sh
go run . --publisher
```

this connects to the message queue named `hello` and sends a `Hello World!` plaintext message
