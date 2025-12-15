# FluxFeed - RSS Aggregator & Reader

A high-performance, concurrent RSS feed aggregator and reader built with Go. FluxFeed efficiently scrapes multiple RSS feeds, stores articles in a PostgreSQL database, and provides a RESTful API for users to manage feeds and read aggregated content.

## Features

- **Multi-Feed Management**: Create and manage multiple RSS feed subscriptions
- **Concurrent Feed Scraping**: Background worker processes (configurable concurrency) that continuously fetch and update feeds
- **User Authentication**: API key-based authentication for secure access
- **Feed Following**: Users can follow/unfollow feeds to customize their reading experience
- **Post Aggregation**: Fetches and stores posts from followed feeds with full-text content
- **Pagination Support**: Retrieve posts with configurable limit and offset
- **Error Handling**: Robust error handling with comprehensive logging
- **Database Persistence**: PostgreSQL database for reliable data storage

## Tech Stack

- **Language**: Go 1.22+
- **Router**: Chi (chi-router)
- **Database**: PostgreSQL with sqlc for type-safe queries
- **Migrations**: Goose
- **Authentication**: API Key (Authorization header)
- **CORS**: Enabled for cross-origin requests

## Project Structure

```
.
├── cmd/                          # Command-line applications
├── internal/
│   ├── auth/                     # Authentication utilities
│   ├── database/                 # Generated sqlc code & database models
│   └── handler/                  # HTTP handlers and middleware
├── rss/
│   ├── rss.go                    # RSS feed parser
│   └── scraper.go                # Background feed scraper
├── sql/
│   ├── queries/                  # SQL query definitions
│   └── schema/                   # Database migrations
├── main.go                       # Application entry point
├── go.mod                        # Go module definition
└── README.md                     # This file
```

## Installation & Setup

### Prerequisites

- Go 1.22 or higher
- PostgreSQL 12+
- Goose (for database migrations)
- sqlc (for code generation)

### Database Setup

1. Create a PostgreSQL database and user:
```sql
CREATE DATABASE rss_aggregator_db;
CREATE USER rss_user WITH PASSWORD 'your_password';
ALTER ROLE rss_user WITH SUPERUSER;
GRANT ALL PRIVILEGES ON DATABASE rss_aggregator_db TO rss_user;
```

2. Set environment variables:
```bash
export DATABASE_URL="postgres://rss_user:your_password@localhost:5432/rss_aggregator_db?sslmode=disable"
export PORT=8080
```

3. Run migrations:
```bash
cd sql/schema
goose postgres $DATABASE_URL up
cd ../..
```

4. Generate code from SQL queries:
```bash
sqlc generate
```

### Building & Running

```bash
# Build the application
go build -o bin/fluxfeed

# Run the application
./bin/fluxfeed
```

The server will start on `http://localhost:8080` and automatically:
- Connect to the PostgreSQL database
- Start the RSS feed scraper (after 10 seconds)
- Listen for incoming HTTP requests

## API Endpoints

### Public Endpoints

#### Create User
```
POST /v1/users/create
Content-Type: application/json

{
  "name": "John Doe"
}

Response (201):
{
  "id": "uuid",
  "created_at": "2025-12-15T10:00:00Z",
  "updated_at": "2025-12-15T10:00:00Z",
  "username": "John Doe",
  "api_key": "your_api_key_here"
}
```

#### Get All Feeds
```
GET /v1/feeds/fetch

Response (200):
[
  {
    "id": "uuid",
    "created_at": "2025-12-15T10:00:00Z",
    "updated_at": "2025-12-15T10:00:00Z",
    "name": "Tech News",
    "url": "https://example.com/feed.xml",
    "user_id": "uuid",
    "last_fetched_at": "2025-12-15T10:30:00Z"
  }
]
```

### Protected Endpoints (Require Authorization Header)

All protected endpoints require:
```
Authorization: ApiKey your_api_key_here
```

#### Get Current User
```
GET /v1/users/fetch

Response (200):
{
  "id": "uuid",
  "created_at": "2025-12-15T10:00:00Z",
  "updated_at": "2025-12-15T10:00:00Z",
  "username": "John Doe",
  "api_key": "your_api_key_here"
}
```

#### Create Feed
```
POST /v1/feeds/create
Content-Type: application/json

{
  "name": "Tech News",
  "url": "https://example.com/feed.xml"
}

Response (201):
{
  "id": "uuid",
  "created_at": "2025-12-15T10:00:00Z",
  "updated_at": "2025-12-15T10:00:00Z",
  "name": "Tech News",
  "url": "https://example.com/feed.xml",
  "user_id": "uuid",
  "last_fetched_at": null
}
```

#### Follow Feed
```
POST /v1/feeds-follow/create
Content-Type: application/json

{
  "feed_id": "uuid"
}

Response (201):
{
  "id": "uuid",
  "created_at": "2025-12-15T10:00:00Z",
  "updated_at": "2025-12-15T10:00:00Z",
  "user_id": "uuid",
  "feed_id": "uuid"
}
```

#### Get User's Feed Follows
```
GET /v1/feeds-follow/fetch

Response (200):
[
  {
    "id": "uuid",
    "created_at": "2025-12-15T10:00:00Z",
    "updated_at": "2025-12-15T10:00:00Z",
    "user_id": "uuid",
    "feed_id": "uuid"
  }
]
```

#### Get User's Feed Posts
```
GET /v1/feeds-follow/user

Response (200):
[
  {
    "id": "uuid",
    "created_at": "2025-12-15T10:00:00Z",
    "updated_at": "2025-12-15T10:00:00Z",
    "title": "Article Title",
    "description": "Article description or content preview",
    "published_at": "2025-12-15T09:00:00Z",
    "url": "https://example.com/article",
    "feed_id": "uuid"
  }
]
```

#### Unfollow Feed
```
DELETE /v1/feeds-follow/delete/{feedFollowID}

Response (200):
{
  "status": "deleted"
}
```

#### Health Check
```
GET /v1/health

Response (200):
{
  "status": "ok"
}
```

## Usage Example

### 1. Create a User
```bash
curl -X POST http://localhost:8080/v1/users/create \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe"}'
```

Save the returned `api_key`.

### 2. Create a Feed
```bash
curl -X POST http://localhost:8080/v1/feeds/create \
  -H "Authorization: ApiKey YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name":"Tech Blog",
    "url":"https://www.wagslane.dev/index.xml"
  }'
```

Save the returned `feed_id`.

### 3. Follow a Feed
```bash
curl -X POST http://localhost:8080/v1/feeds-follow/create \
  -H "Authorization: ApiKey YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"feed_id":"FEED_ID"}'
```

### 4. Wait for Scraper to Run
The background scraper runs every 60 seconds. Wait for it to fetch posts.

### 5. Retrieve Your Posts
```bash
curl -X GET http://localhost:8080/v1/feeds-follow/user \
  -H "Authorization: ApiKey YOUR_API_KEY"
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    username TEXT NOT NULL UNIQUE,
    api_key TEXT NOT NULL UNIQUE
);
```

### Feeds Table
```sql
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES users(id),
    last_fetched_at TIMESTAMP
);
```

### Feeds Follow Table
```sql
CREATE TABLE feeds_follow (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    feed_id UUID NOT NULL REFERENCES feeds(id),
    UNIQUE(user_id, feed_id)
);
```

### Posts Table
```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL UNIQUE,
    description TEXT,
    published_at TIMESTAMP NOT NULL,
    url TEXT NOT NULL UNIQUE,
    feed_id UUID NOT NULL REFERENCES feeds(id)
);
```

## Configuration

Environment variables:

- `DATABASE_URL`: PostgreSQL connection string (required)
  - Format: `postgres://user:password@host:port/database?sslmode=disable`
- `PORT`: HTTP server port (default: 8080)

## Background Scraper

The RSS scraper runs in the background with the following configuration:

- **Concurrency**: 10 goroutines (configurable in `main.go`)
- **Interval**: 60 seconds between scraping runs
- **Startup Delay**: 10 seconds (allows time for database connections)
- **Behavior**: Fetches feeds ordered by `last_fetched_at` (null feeds first)

## Architecture

### Authentication Flow

1. User creates account → receives `api_key`
2. User includes `api_key` in `Authorization: ApiKey <key>` header
3. Middleware validates key against database
4. User object stored in request context
5. Handler retrieves user from context and processes request

### Feed Scraping Flow

1. Ticker triggers every 60 seconds
2. Queries next N feeds to scrape (ordered by oldest fetch first)
3. For each feed, spawns goroutine to:
   - Fetch RSS feed content
   - Parse XML
   - Extract posts
   - Store in database (skips duplicates via unique constraint)
   - Update `last_fetched_at` timestamp

### Post Aggregation

Posts are queried via JOIN between:
- `posts` table (the content)
- `feeds_follow` table (user's subscriptions)
- `feeds` table (metadata)

Results are ordered by `published_at` descending and paginated.

## Error Handling

- **400**: Bad request (invalid JSON, missing fields)
- **403**: Forbidden (invalid API key, authentication failed)
- **404**: Not found (endpoint doesn't exist)
- **500**: Internal server error (database error, scraper error)

All errors return JSON with an `error` field.

## Development

### Running Tests
```bash
go test ./...
```

### Formatting Code
```bash
go fmt ./...
```

### Linting
```bash
go vet ./...
```

## Future Enhancements

- [ ] Search/filter posts by title or content
- [ ] User preferences (sorting, notification settings)
- [ ] Feed categorization/tagging
- [ ] Full-text search on post content
- [ ] WebSocket support for real-time updates
- [ ] User authentication (OAuth2)
- [ ] Rate limiting
- [ ] Caching layer (Redis)
- [ ] Feed health monitoring
- [ ] User feed recommendations

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Commit changes: `git commit -am 'Add your feature'`
3. Push to branch: `git push origin feature/your-feature`
4. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

Built with ❤️ by the Murtaza Patel
