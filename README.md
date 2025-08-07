# Chat Insights Application

This is a simple insights application built with Go and Echo framework.

## Configuration

The application can be configured using a `config.yaml` file in the root directory. An example configuration is shown below:

```yaml
address: ":8082"
api_key: "insights-api-key"
postgres:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "mysecretpassword"
  dbname: "postgres"
  sslmode: "disable"
  schema: "insights"```

