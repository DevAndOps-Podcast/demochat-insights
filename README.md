# Chat Insights Application

This is a simple insights application built with Go and Echo framework.

## Configuration

The application can be configured using a `config.yaml` file in the root directory. An example configuration is shown below:

```yaml
address: ":8081"
api_key: "your-api-key"
```

- `address`: The address and port the server will listen on (e.g., `:8081`).
- `debug`: A boolean indicating whether debug mode is enabled.

## Getting Started

To run the application:

1. **Build the application:**

    ```bash
    go build -o demochat-insights .
    ```

2. **Run the application:**

    ```bash
    ./demochat-insights
    ```

The application will start on the address specified in `config.yaml` (default: `:8081`).
