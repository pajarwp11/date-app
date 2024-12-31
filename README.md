# Date App

Welcome to **Date App**! This application allows users to sign up, like other users, and match with them.

## Structure

This service consist of several folder

### config
Includes all necessary configuration for the service, example: env, mysql, redis

### utils
Includes all necessary code that can be utilised in various logic

### models
Includes all necessary struct model, divided by entity

### handler
This layer will handle incoming request from client including validation and returning response

### service
This layer will handle business logic

### repository
This layer will handle data storage and fetching, for example from database, cache,or external service


## How To Run
1. Create .env file and copy .env.example content and adjust the env according to your own configuration.
2. **go run main.go**