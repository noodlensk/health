# health

TBD

## How it works

For now, it supports only custom health check scripts(commands) defined in the config. Command should
generate `report.json` file with list of all executed checks and their result.

File format:

```json
{
  "Checks": [
    {
      "Error": "got wrong response",
      "Group": "Simple health check",
      "Name": "Simple check scenario A",
      "Priority": "HIGH",
      "Status": "FAILED",
      "Took": "0.009698s"
    },
    {
      "Error": "",
      "Group": "Simple health check",
      "Name": "Simple check scenario B",
      "Priority": "HIGH",
      "Status": "PASSED",
      "Took": "0.000641s"
    },
    {
      "Error": "",
      "Group": "Simple health check",
      "Name": "Simple check scenario C",
      "Priority": "HIGH",
      "Status": "PASSED",
      "Took": "0.000596s"
    }
  ],
  "Created": "2021-03-14T20:55:38.885+01:00",
  "Status": "FAILED",
  "Took": "0.002089s"
}
```

Examples of integrations:

- [Cucumber(Ruby)](./examples/cucumber/README.md)
- [Bash](./examples/bash)

## Development

```bash
make dep
make build-server
```

## TODO

- [ ] Release process (build binary, create GH release, put binary as artefact)
- [ ] OpsGenie close alert
- [ ] Report when can't execute command (OpsGenie, Prometheus)
- [ ] Introduce built-in checks for HTTP, gRPC, TCP, etc
- [ ] Separate server and agents
- [ ] Make server store state
- [ ] API (run test suite, get list of test suites, SLA, etc)
- [ ] Tests

