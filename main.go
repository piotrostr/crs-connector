package main

import (
	"context"
	"fmt"
	"log"

	retail "cloud.google.com/go/retail/apiv2alpha"
	"google.golang.org/api/option"
	retailpb "google.golang.org/genproto/googleapis/cloud/retail/v2alpha"
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

func main() {
	ctx := context.Background()

	productClient := GetProductClient(ctx)

	name := fmt.Sprintf("projects/%s/locations/global/catalogs/default_catalog/branches/default_branch/products/%s", PROJECT_ID, PRODUCT_ID)

	defer productClient.Close()

	product, err := productClient.UpdateProduct(ctx, &retailpb.UpdateProductRequest{
		Product: &retailpb.Product{
			Name:        name,
			Title:       "Test product",
			Description: "Test product description",
			Categories:  []string{"Testing"},
		},
		AllowMissing: true,
		// optional: UpdateMask: &retailpb.FieldMask{
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(product)
}
