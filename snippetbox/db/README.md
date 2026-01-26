# Database Directory

This directory contains database-related files for the snippetbox project.

## Structure

- `init/` - SQL initialization scripts that run when the Docker container is first created
  - `01_schema.sql` - Creates all necessary tables (snippets, users, sessions) and inserts sample data

## Usage

The SQL files in the `init/` directory are automatically executed by MySQL when the Docker container starts for the first time. They are executed in alphabetical order.

## Tables Created

1. **snippets** - Stores code snippets with title, content, and expiration
2. **users** - Stores user accounts with authentication details
3. **sessions** - Stores session data for user authentication
