# Product Service

[<-- Back to Main readme.md File](../readme.md)

The Product Service is responsible for managing products within the order management system.

## GraphQL API Implementation

The Product Service exposes a GraphQL API for interacting with product-related data.

### Create Product

```graphql
mutation {
  createProduct(input: { name: "Smart Thermostat", description: "Energy-efficient smart thermostat that allows you to control your home’s temperature remotely. Features learning capabilities to adjust to your schedule and save energy.", price: 149.99 }) {
    id
    name
    description
    price
  }
}
```

### GetProductByID

```graphql
query {
  getProductByID(id: "d518ff72-05e4-4b18-b81a-7d397d3a5ff2") {
    id
    name
    description
    price
  }
}
```

### ListProducts

```graphql
query {
  listProducts(
    pagination: {limit: 5, offset: 0}
  ) {
    id
    name
    description
    price
  }
}
```

### ListProductsWithIDs

```graphql
query {
  listProductsWithIDs(
    ids: [
      "528c443b-42a2-412a-a460-270c117b805f",
      "b642bee3-07d5-4875-a05d-0ce9021cf345"
    ],
    pagination: {
      limit: 10,
      offset: 0
    }
  ) {
    id
    name
    description
    price
  }
}

```

### SearchProducts

```graphql
query {
  searchProducts(query: "bluetooth", pagination: { limit: 10, offset: 0 }) {
    id
    name
    description
    price
  }
}

```
