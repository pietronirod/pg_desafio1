# Rate Limiter Middleware with Redis

This project implements a Rate Limiter middleware using Redis for controlling request rates based on IP addresses or access tokens. It is designed to help protect your application from abuse by limiting the number of requests a user can make in a given time frame.

## Project Structure

```
rate-limiter
├── cmd
│   └── main.go          # Entry point of the application
├── internal
│   ├── middleware
│   │   └── rate_limiter.go  # Rate limiter middleware implementation
│   ├── redis
│   │   └── client.go    # Redis client for handling connections
│   └── config
│       └── config.go    # Configuration settings loader
├── go.mod                # Module definition and dependencies
└── README.md             # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.16 or later
- Redis server

### Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd rate-limiter
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

### Configuration

The application uses environment variables or a `.env` file for configuration. You can set the following variables:

- `REQUEST_LIMIT`: Maximum number of requests allowed within the specified time frame.
- `BLOCK_DURATION`: Duration for which the IP or token will be blocked after exceeding the limit.

### Running the Application

To run the application, execute the following command:

```
go run cmd/main.go
```

### Usage

Once the server is running, it will automatically apply the rate limiting middleware to incoming requests based on the configured limits.

### Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

### License

This project is licensed under the MIT License. See the LICENSE file for details.