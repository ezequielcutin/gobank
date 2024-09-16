# GoBank

GoBank is a bank account management system built from scratch using Go. It provides a RESTful API for managing user accounts, including functionalities for creating, retrieving, deleting accounts, and transferring funds. The project is designed to be lightweight, utilizing minimal external libraries and is containerized using Docker.

## Features

- **Create Account**: Allows users to create a new bank account.
- **Get Accounts**: Retrieves a list of all bank accounts.
- **Get Account by ID**: Fetches details of a specific account using its ID.
- **Delete Account**: Deletes an existing bank account.
- **Transfer Funds**: Facilitates transferring money between accounts.
- **Overdraft Protection**: Prevent users from withdrawing more funds than they have.
- **Monopoly Bank Account**: A special account for testing purposes with a balance of 999,999,999.

### Overdraft Protection
To ensure that users cannot withdraw more funds than they have, the application implements overdraft protection. When a transfer is initiated, the system checks the sender's account balance before proceeding. If the balance is insufficient, the transfer will be rejected, and an error message will be returned. This feature helps maintain the integrity of user accounts and prevents negative balances.

### Monopoly Bank Account
The application includes a special account known as the "Monopoly Bank," which is designed for testing purposes. This account has a balance of **999,999,999** and serves as a source of funds for transfers. 

- **Account ID**: If running `make run` in the terminal for the first time, it will be ID of 1. Modify as you wish in `main.go` to avoid duplicates.
- **Balance**: 999,999,999
- **Purpose**: To facilitate testing of the transfer functionality without the risk of running out of funds.

The system checks for the existence of the Monopoly Bank account upon startup. If it does not exist, the account will be created automatically. This prevents duplicate accounts from being created during multiple runs of the application.

## Technologies Used

- Go (version 1.22.5)
- Gorilla Mux for routing (indirectly referenced)
- PostgreSQL for data storage
- Docker for containerization

## Getting Started

### Prerequisites

- **Node.js** (v14 or above)
- **npm** or **yarn**
- **Go** (v1.16 or above)
- **PostgreSQL** database


### Installation

1. Clone the repository:
   ```
   git clone https://github.com/ezequielcutin/gobank.git
   ```

2. Navigate to the project directory:
   ```
   cd gobank
   ```

3. Install Backend dependencies:
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

To start the API server, simply run in the terminal (in the root directory):
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

## Frontend Requirements

1. Navigate to the frontend directory:
   ```
   cd frontend
   ```

2. Install frontend dependencies:
   ```
   npm install
   ```

3. Start the frontend application:
   ```
   npm start
   ```

   Open [http://localhost:3001](http://localhost:3001) (or whichever port the frontend is running on) to view it in your browser.
   Make sure the backend is running before running the frontend, otherwise they will not be able to communicate.
   You can run the frontend and backend in different terminals.


## API Endpoints

### User Authentication

- **Register**

  ```
  POST /register
  ```

  **Request Body:**

  ```json
  {
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "password": "securepassword"
  }
  ```

- **Login**

  ```
  POST /login
  ```

  **Request Body:**

  ```json
  {
    "email": "john.doe@example.com",
    "password": "securepassword"
  }
  ```

### Balance

- **Get Balance**

  ```
  GET /balance/{userId}
  ```

  **Response:**

  ```json
  {
    "balance": 1000
  }
  ```

### Transactions

- **Get Transactions**

  ```
  GET /transactions/{userId}
  ```

  **Response:**

  ```json
  [
    {
      "id": 1,
      "amount": -100,
      "type": "Sent",
      "createdAt": "2023-10-01T12:34:56Z",
      "otherUserId": 2
    },
    {
      "id": 2,
      "amount": 200,
      "type": "Received",
      "createdAt": "2023-10-02T09:21:43Z",
      "otherUserId": 3
    }
  ]
  ```

### Users

- **Get Users**

  ```
  GET /account
  ```

- **Get User by Email**

  ```
  GET /user-by-email/{email}
  ```

- **Get User Details**

  ```
  GET /user-details/{email}
  ```


## Frontend Components

### Dashboard

- **Description**: The main dashboard that aggregates all user-specific components such as balance, transfer funds, available users, and transaction history.

### CurrentBalanceCard

- **Description**: Displays the user's current balance with a dynamic gradient background that responds to mouse movements.
- **Props**:
  - `refreshTrigger`: Triggers data refresh when updated.

### TransferFundsCard

- **Description**: Enables users to transfer funds to other users.
- **Props**:
  - `onTransferSuccess`: Callback function to refresh data upon successful transfer.

### TransactionHistoryCard

- **Description**: Shows a history of all transactions, indicating whether funds were sent or received, along with corresponding user IDs.
- **Props**:
  - `refreshTrigger`: Triggers data refresh when updated.
  - `userId`: The ID of the current user.

### AvailableUsersCard

- **Description**: Lists all users available for transactions.
- **Props**:
  - `currentUserId`: The ID of the current user to filter out from the list.

### UserList

- **Description**: Displays a list of all user accounts with options to delete accounts.
- **Props**:
  - `onDelete`: Callback function to refresh data upon account deletion.


## Contributing

Contributions are welcome! Please follow these steps:

1. **Fork the Repository**

2. **Create a Feature Branch**

   ```bash
   git checkout -b feature/YourFeatureName
   ```

3. **Commit Your Changes**

   ```bash
   git commit -m "Add some feature"
   ```

4. **Push to the Branch**

   ```bash
   git push origin feature/YourFeatureName
   ```

5. **Open a Pull Request**


## License

This project is licensed under the MIT License.
