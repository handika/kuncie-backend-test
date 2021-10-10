# Kuncie Backend Take Home Test

Have you shopped online? Let’s imagine that you need to build the checkout backend service that will support different promotions with the given inventory.

Build a checkout system with these items:

![gambar](https://user-images.githubusercontent.com/1314588/136423748-a7aa28b6-2d10-4f06-b604-f0fe011a2678.png)

**The system should have the following promotions:**
- Each sale of a MacBook Pro comes with a free Raspberry Pi B
- Buy 3 Google Homes for the price of 2
- Buying more than 3 Alexa Speakers will have a 10% discount on all Alexa speakers

**Example Scenarios:**
- Scanned Items: MacBook Pro, Raspberry Pi B
  - Total: $5,399.99
- Scanned Items: Google Home, Google Home, Google Home
  - Total: $99.98
- Scanned Items: Alexa Speaker, Alexa Speaker, Alexa Speaker
  - Total: $295.65

Please write it in Golang or Node with a CI script that runs tests and produces a binary.

Finally, imagine that adding items to cart and checking out was a backend API. Please design a schema file for GraphQL on how you would do this.

Thank you for your time and we look forward to reviewing your solution. If you have any questions, please feel free to contact us. Please send us a link to your git repo.

# Database Schema

![schema](https://user-images.githubusercontent.com/1314588/136503808-c7c479d4-2122-4fc6-9fd1-8327a9355ccf.png)

# Building and Running The App

**Prerequisites:**

    Go 1.15.2
    Docker
    Docker Compose
    Golang migrate (https://github.com/golang-migrate/migrate)

**Step 1 Checkout**

```
$ git clone https://github.com/handika/kuncie-backend-test.git
$ cd kuncie-backend-test
```

**Step 2 Start MySQL Service**

```
$ docker-compose up
```

**Step 3 Run Migration**

```
$ migrate -database mysql://backend:backend@/backend -path ./mysql up
```
**Step 4 GraphQL Playground**

```
http://localhost:8080/
```

**Step 5 Calling APIs**

Create Transaction

```
curl 'http://localhost:8080/query' -H 'Accept-Encoding: gzip, deflate, br' -H 'Content-Type: application/json' -H 'Accept: application/json' -H 'Connection: keep-alive' -H 'DNT: 1' -H 'Origin: http://localhost:8080' --data-binary '{"query":"# Write your query or mutation here\nmutation {\n  createTransaction(\n    input: {\n      userId: 1\n      details: [{\n        productId: 3\n        qty: 3\n      },\n      {\n        productId: 1\n        qty: 3\n      }]\n    }\n  ) {\n    id\n    userId\n    grandTotal\n    details {\n      productId\n      price\n      qty\n      subTotal\n      discount\n    }\n  }\n}\n"}' --compressed
```

```
mutation {
  createTransaction(
    input: {
      userId: 1
      details: [{
        productId: 3
        qty: 3
      },
      {
        productId: 1
        qty: 3
      }]
    }
  ) {
    id
    userId
    grandTotal
    details {
      productId
      price
      qty
      subTotal
      discount
    }
  }
}

```

Get Transaction by ID

```
curl 'http://localhost:8080/query' -H 'Accept-Encoding: gzip, deflate, br' -H 'Content-Type: application/json' -H 'Accept: application/json' -H 'Connection: keep-alive' -H 'DNT: 1' -H 'Origin: http://localhost:8080' --data-binary '{"query":"# Write your query or mutation here\nquery {\n  transactionByID(\n    id: 1\n  ) {\n    id\n    grandTotal\n    userId\n    details {\n      productId\n      price\n      qty\n      subTotal\n      discount\n    }\n  }\n}\n"}' --compressed
```
```
query {
  transactionByID(
    id: 31
  ) {
    id
    grandTotal
    userId
    details {
      productId
      price
      qty
      subTotal
      discount
    }
  }
}

```
