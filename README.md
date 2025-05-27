![eyeOne Logo](logo.png)
# eyeOne


**eyeOne** is a lightweight trading API built with Go and Gin. It provides endpoints for creating and managing orders, retrieving balances, and accessing order books. This project is designed for developers seeking a simple and extensible trading backend.

---

## 🚀 Features

- **Create Orders**: Place new buy or sell orders with specified parameters.
- **Cancel Orders**: Remove existing orders using their unique identifiers.
- **Retrieve Balances**: Check the balance of specific assets.
- **Access Order Books**: View the current order book for trading pairs.
- **RESTful API**: Built with Gin for efficient HTTP request handling.

---

## 📦 Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/erfuuan/eyeOne.git
   cd eyeOne
   ```

2. **Set Up Environment Variables**:

   Create a `.env` file in the root directory and configure your environment variables as needed.

3. **Install Dependencies**:

   ```bash
   go mod tidy
   ```

4. **Run the Application**:

   ```bash
   go run cmd/main.go
   ```

   The server will start on `http://localhost:8080`.

---
## 🐳 Deployment with Docker

To deploy **eyeOne** using Docker, follow these steps:

### Build the Docker image:

```bash
docker build -f build/Dockerfile -t eyeone .
```

### Run the Docker container:
```bash
docker run -d -p 3000:3000 --name eyeOne eyeOne:latest
```
---
## 📖 API Endpoints
### 1. Create Order
- **Description:** Place a new order.
- **Response:** Returns the order ID upon successful creation.
---

### 2. Cancel Order

- **Endpoint:** `DELETE /order-book/:id`
- **Description:** Cancel an existing order by its ID.
- **Parameters:**
  - `id`: The unique identifier of the order to cancel.
- **Response:** Confirmation of order cancellation.

---
### 3. Get Balance
- **Description:** Retrieve the balance for a specific asset.
- **Response:** Returns the balance amount for the specified asset.

---
### 4. Get Order Book
- **Description:** Fetch the order book for a trading pair.
- **Response:** Returns the current order book data for the specified symbol.

---

## 🧪 Testing with Postman

To facilitate testing, you can use the following Postman collection:

1. **Import the Collection:**

   Download the `eyeOne.postman_collection.json` file from the repository.

2. **Open Postman:**

   Launch Postman and click on **Import**.

3. **Import the File:**

   Select the downloaded `eyeOne.postman_collection.json` file to import the collection.

4. **Use the Endpoints:**

   The collection includes pre-configured requests for all API endpoints. Modify the request parameters as needed and send the requests to test the API.

---