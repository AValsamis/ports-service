# Ports Service Documentation

## Overview
The Ports Service is responsible for managing ports based on the data provided in a ports.json file. It can create new records in a database or update existing ones, ensuring that the resulting storage represents the latest version found in the JSON file. The service operates with limited resources, such as 200MB RAM, and handles files of unknown size, which can contain several million records.

## Implementation Details

### Domain Model
The domain model represents the core business entities and logic of the port service. In this implementation, the domain model is directly used for the repository as well for simplicity.
> In a production project it's recommended to introduce separate models for the repository layer to maintain separation of concerns and flexibility in future modifications.

### Repository
The repository is responsible for persisting and retrieving ports. It provides methods for creating new records and updating existing ones. The repository implementation uses a Redis database to store the ports.

Redis as a storage solution provides several advantages over an in-memory solution:
1. Persistence: Redis allows data to be persisted to disk, ensuring that the port records are not lost in case of service restarts or failures.
2. Scalability: Redis is designed to handle large datasets efficiently and can scale horizontally to support increasing port records.
3. Concurrency: Redis handles concurrent access to data, enabling multiple clients to read and write simultaneously.

### Service
The service layer acts as an interface between the repository and the domain logic. It provides methods for loading ports from the ports.json file, decoding the JSON data, and upserting the ports into the repository. It also includes validation to ensure the integrity of the port data.

## Running the Application

### Prerequisites
- Go 1.19 or later installed 
- Docker

### Instructions
1. Clone the repository and navigate to the project directory.
2. Run the following command to build the application:
   ```shell
   make docker-build
   ```

3. To run the application using Docker, use the following commands:
   ```shell
   make docker-run
   ```

   This will build the Docker image and run the application inside a container, loading the ports from the ports.json file.

4. To clean up the containers make sure to run the following command:
   ```shell
   make docker-down
   ```
   
## Testing
The application includes unit and integrations tests to verify the functionality of the different components. To run the tests, use the following command:
```shell
make test
```

The tests include coverage analysis and utilize the race detector to identify potential data race conditions.
> In a production project we would also need end-to-end tests, performance tests etc.

## Linting and Formatting
The project includes linting and formatting tools to maintain code quality and consistency. To run the linters and format the code, use the following commands:
```shell
make lint
make fmt
```

These commands will check for common issues, apply formatting, and ensure adherence to coding standards.
> In a production project we would add these code quality measures in a CI pipeline to ensure that they are not missed and whatever is committed adheres to those coding standards.

## Additional Notes
- The application is designed to gracefully handle termination signals to ensure proper shutdown.
- The Redis database used for storage provides persistence, scalability, concurrency support, and efficient memory management, making it suitable for managing large port datasets.