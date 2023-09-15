package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	retail "cloud.google.com/go/retail/apiv2alpha"
	"google.golang.org/api/option"
	retailpb "google.golang.org/genproto/googleapis/cloud/retail/v2alpha"

	amqp "github.com/rabbitmq/amqp091-go"
)

const PROJECT_ID = "piotrostr-resources"
const PRODUCT_ID = "test-product"

func GetProductClient(ctx context.Context) *retail.ProductClient {
	opts := []option.ClientOption{
		option.WithCredentialsFile("service-account.json"),
	}
	client, err := retail.NewProductClient(ctx, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func ConvertToRetail(product *any) *retailpb.Product {
	return nil
}

func Publish(ctx context.Context) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	body := "Hello World!"
	err = ch.PublishWithContext(
		ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(body),
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func Subscribe(ctx context.Context, productClient *retail.ProductClient) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Waiting for messages...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Printf("Received a message: %s\n", d.Body)
			fmt.Printf("Changing the test %s title\n", PRODUCT_ID)
			product, err := productClient.UpdateProduct(ctx, &retailpb.UpdateProductRequest{
				Product: &retailpb.Product{
					Name: fmt.Sprintf(
						"projects/%s/locations/global/catalogs/default_catalog/branches/default_branch/products/%s",
						PROJECT_ID,
						PRODUCT_ID,
					),
					Title:       string(d.Body),
					Description: "Test product description",
					Categories:  []string{"Testing"},
				},
				AllowMissing: true,
				// optional: UpdateMask: &retailpb.FieldMask{}
			})
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Product updated: %s\n", product)
		}
	}()

	<-forever
}

func main() {
	publisher := flag.Bool("publisher", false, "publish message")
	subscriber := flag.Bool("subscriber", false, "receive messages")
	flag.Parse()

	ctx := context.Background()

	productClient := GetProductClient(ctx)
	defer productClient.Close()

	if *publisher {
		fmt.Println("Publishing message...")
		Publish(ctx)
	} else if *subscriber {
		fmt.Println("Subscribing & updating...")
		Subscribe(ctx, productClient)
	}
}
