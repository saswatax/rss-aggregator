# RSS Aggregator

An RSS feed aggregator built with Go that allows users to create and subscribe to RSS feeds, featuring account management and JWT-based authentication with automatic feed updates.

## 🚀 Features

- Subscribe to RSS feeds.
- Fetch and store RSS feed data.
- Manage user accounts and API keys.
- Follow feeds and track updates.
- Authentication for secure access.

## 🛠️ Technologies Used

- [Golang](https://golang.org)
- [PostgreSQL](https://www.postgresql.org)
- [Chi Router](https://github.com/go-chi/chi)
- [JWT Authentication](https://jwt.io)

## 🛠️ Setup and Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/your-username/rss-aggregator.git
   cd rss-aggregator
   ```

2. Install dependencies:

   ```sh
   go mod download
   ```

3. Build the application:

   ```sh
   go build -o ./bin/ ./cmd/api
   ```

4. Configure environment variables:
   Create a `.env` file in the root directory with the following:

   ```
   PORT=8080
   DB_URL=postgresql://username:password@localhost:5432/rss_aggregator?sslmode=disable
   ```

5. Run the application:
   ```sh
   ./bin/api
   ```
