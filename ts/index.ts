import { protos, ProductServiceClient } from "@google-cloud/retail";

// read JSON product

// map schema to Product proto
const product = new protos.google.cloud.retail.v2.Product({
  title: "my-product",
  type: "PRIMARY",
  categories: ["my-category"],
  attributes: {
    "my-attribute-key": {
      text: ["my-attribute-value"],
    },
    // ...
  },
});

const client = new ProductServiceClient();

// send the request with the product from the queue
async function createProduct() {
  client.updateProduct({
    product: product,
    allowMissing: true,
  });
}
