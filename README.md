# gometrics

# usage example

```go
var metrics gometrics.Metrics
metrics.Init("", "metrics_config.yaml")
```

# yaml configuration example
```yaml
port: 8000
processing_total:
  - type: "process"
    process: "test"
    trigger_time: "1m"
    trigger_count: "1"
error_total:
  - type: "error"
    process: "test"
    trigger_time: "1m"
    trigger_count: "1"
warning_total:
  - type: "warn"
    process: "test"
    trigger_time: "1m"
    trigger_count: "1"
```