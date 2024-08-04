# k6 Load Testing Project

This project is set up to perform load testing using k6.

## Directory Structure

- `scripts/`: Contains k6 test scripts.
- `results/`: Directory to store test results.
- `Dockerfile`: Docker configuration for running k6 tests.
- `.env`: Environment variables.
- `config.json`: Configuration file for k6.
- `README.md`: Project documentation.

## Running the Tests

1. Build the Docker image:

    ```sh
    docker build -t k6-load-test .
    ```

2. Run the k6 test:

    ```sh
    docker run --rm --env-file .env k6-load-test
    ```

3. (Optional) Save the results:

    ```sh
    docker run --rm --env-file .env -v $(pwd)/results:/results k6-load-test run /scripts/get-users.js --out json=/results/result.json
    ```

## Configuration

- Modify the `scripts/` files to add or update test scenarios.
- Update `.env` with your environment variables.
- Adjust `config.json` for different test settings.
