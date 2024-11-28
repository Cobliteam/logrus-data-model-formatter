# Logrus data model formatter
A Logrus formatter for structuring logs according to Cobli's [standardized data model](https://www.notion.so/cobli/Modelo-de-dados-normaliza-o-de-logs-11f682738c9c8078aa00f4edc69dcd6b). This formatter ensures consistency and normalization of log entries across applications, making log analysis easier and more efficient.


## Installation
```bash
go get github.com/Cobliteam/logrus-data-model-formatter
```

## Usage

- Configure formatter and DD hook
    ```go
    package configuration

    import (
        data_model "github.com/Cobliteam/logrus-data-model-formatter"
        "github.com/sirupsen/logrus"
        dd_logrus "gopkg.in/DataDog/dd-trace-go.v1/contrib/sirupsen/logrus"
    )

    type DDLogger struct {
    }

    func (l DDLogger) Log(msg string) {
        logrus.Debug(msg)
    }

    func SetupLogger(logLevel string) {
        logrus.SetFormatter(&data_model.LogrusDataModelFormatter{})
        logrus.AddHook(&dd_logrus.DDContextLogHook{})
    }
    ```
- Initialize logger and DD tracer with logger
    ```go
    package main

    import (
        // ...
    )

    func main() {
        // ...
        config := configuration.NewApiConfig()
        configuration.SetupLogger(config.Server.LogLevel)
        if config.Server.TraceEnabled {
            tracer.Start(
                tracer.WithLogger(&configuration.DDLogger{}),
            )
            defer tracer.Stop()
        }
        // ...
    }
    ```
- Log away
    ```go
    log.WithContext(stream.Context()).WithFields(logrus.Fields{
    "thing_id": in.GetThingId(),
    }).Infof("Received location stream request for thingID=<%s>, startTime=<%d>, endTime=<%d>",
        in.GetThingId(), in.GetStartTimestamp(), in.GetEndTimestamp())
    ```
    - Output
        ```bash
        {"custom":{"thing_id":"thing_id_1"},"dd":{"span_id":4256710396338970337,"trace_id":4256710396338970337},"level":"info","message":"Received location stream request for thingID=<thing_id_1>, startTime=<100>, endTime=<101>","timestamp":"2024-11-27T17:39:20-03:00"}
        ```
