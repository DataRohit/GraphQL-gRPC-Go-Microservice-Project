# Account Service

The Account Service is responsible for managing user accounts within the order management system.

## GraphQL API Implementation

The Account Service exposes a GraphQL API for interacting with account-related data.

### Create Account

```graphql
mutation {
  createAccount(
    name: "Rohit Ingole"
    email: "datarohit@example.com"
  ) {
    id
    name
    email
    createdAt
    updatedAt
  }
}
```

### Get Account by ID

```graphql
query {
  getAccountByID(
    id: "d88ff73c-7563-42aa-896e-f20ed09c1f30"
  ) {
    id
    name
    email
    createdAt
    updatedAt
  }
}
```

### Get Account by Email

```graphql
query {
  getAccountByEmail(
    email: "datarohit@example.com"
  ) {
    id
    name
    email
    createdAt
    updatedAt
  }
}
```

### List Accounts

```graphql
query {
  listAccounts(
    pagination: {limit: 10, offset: 0}
  ) {
    id
    name
    email
    createdAt
    updatedAt
  }
}
```
