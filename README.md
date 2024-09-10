# GoBank

GoBank is a bank account management system built from scratch using Go. It provides a RESTful API for managing user accounts, including functionalities for creating, retrieving, deleting accounts, and transferring funds. The project is designed to be lightweight, utilizing minimal external libraries and is containerized using Docker.

## Features

- **Create Account**: Allows users to create a new bank account.
- **Get Accounts**: Retrieves a list of all bank accounts.
- **Get Account by ID**: Fetches details of a specific account using its ID.
- **Delete Account**: Deletes an existing bank account.
- **Transfer Funds**: Facilitates transferring money between accounts.

## Technologies Used

- Go (version 1.22.5)
- Gorilla Mux for routing (indirectly referenced)
- PostgreSQL for data storage
- Docker for containerization

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/ezequielcutin/gobank.git
   ```

2. Navigate to the project directory:
   ```
   cd gobank
   ```

3. Install dependencies:
   ```
   go mod tidy
   ```

## Docker Setup

Before running the application, set up the PostgreSQL database using Docker. Run the following command in your terminal:
    ```
    docker run --name some-postgres -e POSTGRES_PASSWORD=gobank -p 5432:5432 -d postgres
    ```


Make sure the password you set in the terminal matches the password in the connection string in `storage.go`:
    
    
    func newPostGresStore() (PostgresStore, error) {
    connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
    


## Usage

To start the API server, simply run in the terminal:
    ```
    make run
    ```


The server will run on the specified port, and you can access the API endpoints using tools like Postman or curl.

## API Endpoints

- `GET /account`: Retrieve all accounts.
- `GET /account/{id}`: Retrieve a specific account by ID.
- `POST /account`: Create a new account.
- `DELETE /account/{id}`: Delete an account by ID.
- `POST /transfer`: Transfer funds between accounts (implementation to be added).


## Testing with Postman

You can use Postman to test the API endpoints. Import the provided Postman collection (if available) or manually create requests to interact with the API.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or features.

## License

This project is licensed under the MIT License.