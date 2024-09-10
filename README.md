# Golang Blockchain Application

This is a simple blockchain application implemented in Go. It simulates a library book checkout system using blockchain technology.

## Features

- Create new books
- Add book checkout records to the blockchain
- View the entire blockchain

## Prerequisites

- Go (version 1.15 or later)
- Postman (for testing API endpoints)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/golang-blockchain.git
   cd golang-blockchain
   ```

2. Install dependencies:
   ```
   go mod init golang-blockchain
   go get github.com/gorilla/mux
   ```

3. Build and run the application:
   ```
   go build
   ./golang-blockchain
   ```

The server will start running on `http://localhost:8080`.

## API Endpoints

### 1. Get Blockchain

- **URL:** `/`
- **Method:** GET
- **Description:** Retrieves the entire blockchain.

### 2. Add Block (Book Checkout)

- **URL:** `/`
- **Method:** POST
- **Description:** Adds a new book checkout record to the blockchain.
- **Body:**
  ```json
  {
    "book_id": "123456",
    "user": "John Doe",
    "checkout_date": "2024-09-10"
  }
  ```

### 3. Create New Book

- **URL:** `/new`
- **Method:** POST
- **Description:** Creates a new book record.
- **Body:**
  ```json
  {
    "title": "The Go Programming Language",
    "author": "Alan A. A. Donovan",
    "publish_date": "2015-10-26",
    "isbn": "9780134190440"
  }
  ```

## Using Postman

1. **Get Blockchain:**
   - Open Postman and create a new GET request.
   - Enter the URL: `http://localhost:8080/`
   - Click "Send" to see the entire blockchain.

2. **Add Block (Book Checkout):**
   - Create a new POST request.
   - Enter the URL: `http://localhost:8080/`
   - Go to the "Body" tab, select "raw" and choose "JSON" from the dropdown.
   - Enter the JSON body as shown in the API Endpoints section.
   - Click "Send" to add a new block to the blockchain.

3. **Create New Book:**
   - Create a new POST request.
   - Enter the URL: `http://localhost:8080/new`
   - Go to the "Body" tab, select "raw" and choose "JSON" from the dropdown.
   - Enter the JSON body as shown in the API Endpoints section.
   - Click "Send" to create a new book record.

## Testing the Application

1. Start by creating a new book using the "Create New Book" endpoint.
2. Use the book ID returned from step 1 to create a new checkout record using the "Add Block" endpoint.
3. View the blockchain using the "Get Blockchain" endpoint to see your new records.

## Notes

- This is a simple implementation for educational purposes and is not suitable for production use.
- The blockchain is stored in memory and will reset when the application is restarted.
- Error handling and input validation are minimal in this example.

## Contributing

Feel free to fork this repository and submit pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
