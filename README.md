# Build Your Own Redis


This is my implementation of the ["Build Your Own Redis" ](https://codecrafters.io/challenges/redis) challenge from Codecrafters. I've built a Redis-compatible server using Go, covering the following key features:

## Command Handling
Implemented the basic Redis command protocol to handle common commands like `SET`, `GET`, `DEL`, etc. The server can parse client requests and execute the corresponding operations.

## RDB Persistence
Added support for reading and writing Redis' RDB (Redis Database) files, allowing data to persist across server restarts.

## Replication
Implemented a basic replication mechanism, where a Redis server can function as a replica, replicating data from a master server.

## Streams
Included support for Redis Streams, a data structure that enables advanced message queue and event streaming functionality.

## Transactions
Added transaction support, allowing clients to execute multiple commands as an atomic unit, with rollback capabilities.

This project was a great learning experience, diving deep into the inner workings of Redis and building a functional in-memory data store from scratch.

Feel free to explore the code and let me know if you have any questions or feedback!
