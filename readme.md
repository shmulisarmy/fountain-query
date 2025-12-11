# Fountain Query
When you make a sql query, you'll have a live view into the database with that query, our sql engine when done will be able to handle nested queries, joins and many other query capabilities that are present in sql

## How Data Stays in Sync (High-Level Overview)

The core innovation of this project is ensuring **all clients always see identical data**, regardless of when they connect or in what order updates occur. This is achieved through a sophisticated real-time synchronization system.

### The Synchronization Flow

1. **SQL Query as Observable Source**: Your SQL query becomes an observable that watches for changes in underlying database tables.

2. **EventEmitterTree Tracking**: The backend's `EventEmitterTree` subscribes to these observables and tracks all changes to query results. When data changes (adds, removes, updates), it detects what changed.

3. **WebSocket Broadcasting**: Each change is serialized into a sync message and broadcast to all connected clients via WebSocket. Sync messages include:
   - **Type**: `add`, `remove`, `update`, or `load`
   - **Path**: Hierarchical path in the data structure (e.g., `/person_123/todo/todo_456`)
   - **Data**: JSON-serialized row data
   - **Timestamp**: For ordering guarantees

4. **Client-Side Reconciliation**: Each client receives these sync messages and applies them to reconstruct the exact same data structure. The path-based updates ensure that even nested queries maintain consistency.

### Key Guarantees

- **Order Independence**: Clients can connect at different times and still end up with identical data
- **Idempotent Operations**: Updates can be applied in different orders with the same final result
- **Hierarchical Consistency**: Nested subqueries are maintained through path-based updates
- **Real-Time Updates**: All clients see changes instantly as they occur

### Example: Multi-Client Scenario

```
Client A connects → Receives initial Load message with full snapshot
Client B connects → Receives initial Load message with full snapshot  
(Client A and B now have identical data)

User adds new person → HTTP request to backend
Backend inserts row → SQL query re-evaluates
EventEmitterTree detects change → Broadcasts "add" message
Both Client A and Client B receive the "add" message
Both clients apply update at path /person_123
(Both clients still have identical data)
```

The frontend demo (`frontend/`) provides a live demonstration of this synchronization. You can open multiple browser windows and watch as all windows stay perfectly in sync as you add, remove, or modify data.

## Client SDKs

The project includes multiple client SDKs for consuming live query results in different environments:

- **Go SDK** (`live_db_sdks/go/`) - Thread-safe client with `map[string]any` for Go applications
- **Solid.js SDK** (`live_db_sdks/solid/`) - Fine-grained reactive updates using `createMutable`
- **React SDK** (`live_db_sdks/react/`) - Hooks and class-based APIs with immutable state updates

All SDKs connect via WebSocket and automatically sync with the backend's `eventEmitterTree`, maintaining hierarchical data structures that mirror your SQL queries with nested subqueries.

See [`live_db_sdks/README.md`](live_db_sdks/README.md) for detailed documentation and usage examples.

## Trying to learn but don't where to start?
look no further than the pub_sub directory as thats where we build the core primitive components that when assembled create a great programming world of emerging behaviors.



## How to contribute?
There are multiple tickets available in https://github.com/shmulisarmy/sql-compiler/blob/main/TODO.md

## Have an idea? 
Feel free to share here as we are very open to new ideas


### run with
```bash 
go run .
```

### test with
```bash 
go test ./...
```