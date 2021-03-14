# cucumber

Example of usage running health checks based on cucumber(Ruby).

## How it works

- [CustomFormatter](./features/support/formatter.rb) for proper logs format (one line per scenario start/finish)
- [Reporter](./features/support/reporter.rb) for collecting checks execution data and saving it into `report.json` or
  location defined in `$REPORT_FILE` environment variable
- [Rescue code](./features/step_definitions/simple_checks_steps.rb) which wraps every assertion and collect errors
  into `@collected_errors` array
  

## TODO

- [ ] Get rid of [Rescue code](./features/step_definitions/simple_checks_steps.rb) and do it automatically
- [ ] Define severity of check via tag
