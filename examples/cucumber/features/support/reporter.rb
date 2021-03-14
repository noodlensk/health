# frozen_string_literal: true

Before do |_scenario|
  @collected_errors = []
end

After do |_scenario|
  raise @collected_errors.to_s unless @collected_errors.empty?
end

Before do |_scenario|
  @scenario_started = Time.now
end

After do |scenario|
  # HACK: for obtaining feature name, see here https://github.com/cucumber/cucumber-ruby/issues/1432
  string = File.read(scenario.location.file)
  document = ::Gherkin::Parser.new.parse(string)
  feature = document[:feature]

  CucuGlobals.instance.report[:Checks].push({
                                              Name: scenario.name,
                                              Group: feature.name,
                                              Status: scenario.failed? || !@collected_errors.empty? ? 'FAILED' : 'PASSED',
                                              Priority: 'HIGH', # TODO: make it configurable via tags
                                              Error: @collected_errors.join(', '),
                                              Took: "#{Time.now - @scenario_started}s"
                                            })
end

at_exit do
  report = CucuGlobals.instance.report
  report[:Took] = "#{Time.now - CucuGlobals.instance.report[:Created]}s"

  report[:Status] = 'PASSED' if report[:Checks].filter { |c| c[:Status] == 'FAILED' }.empty?

  File.write(ENV['REPORT_FILE'] || 'report.json', report.to_json)
end
