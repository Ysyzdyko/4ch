# 4ch


Project Overview

1337b04rd is an anonymous imageboard inspired by early internet forums. This project implements a simplified version of modern imageboards with posts, comments, and image uploads.
Key Features

    Anonymous posting (no registration required)

    Threaded comments with replies

    Image uploads via S3-compatible storage

    Unique Rick and Morty avatars assigned per session

    Automatic post archiving based on activity

    Hexagonal architecture design

Technical Stack

    Backend: Go (standard library + PostgreSQL driver)

    Database: PostgreSQL

    Storage: S3-compatible (using triple-s implementation)

    Frontend: Provided HTML templates

    Authentication: Session cookies

    Logging: log/slog package

Setup Instructions
Prerequisites

    Go 1.21+

    PostgreSQL 15+

    MinIO/S3-compatible storage

    Docker (optional)

Installation

    Clone the repository:
    bash

git clone https://github.com/yourusername/1337b04rd.git
cd 1337b04rd

Set up environment variables (create .env file):
env

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=1337b04rd
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET_POSTS=posts
S3_BUCKET_AVATARS=avatars
S3_USE_SSL=false

Build and run:
bash

    go build -o 1337b04rd ./cmd/1337b04rd
    ./1337b04rd --port 8080

Docker Setup (Alternative)
bash

docker compose up

API Endpoints
Posts

    GET / - Catalog page

    GET /archive - Archive page

    GET /post/{id} - View post

    GET /archive/post/{id} - View archived post

    GET /create-post - Create post form

    POST /submit-post - Submit new post

    POST /post/{id}/submit-comment - Add comment

Session Management

    Automatic cookie-based session handling

    Avatar assignment via Rick and Morty API

Database Schema

Key tables:

    posts - Threads with metadata

    comments - Replies to posts

    sessions - User session tracking

    archived_posts - Archived threads

Architecture

Application Core (Domain)
├── Ports
│   ├── PostRepository
│   ├── CommentRepository
│   ├── SessionManager
│   └── ImageStorage
└── Adapters
    ├── PostgreSQLAdapter
    ├── S3StorageAdapter
    ├── RickAndMortyAvatarAdapter
    └── HTTPHandlerAdapter

Testing

Run tests with:
bash

go test -v ./...

Test coverage target: ≥20%
Deployment

    Set up PostgreSQL database

    Configure S3-compatible storage

    Build and run binary:
    bash

    ./1337b04rd --port 80

Usage

    Visit / to view active posts

    Click "Create Post" to start a new thread

    Click any post to view and comment

    Archived posts are available at /archive
