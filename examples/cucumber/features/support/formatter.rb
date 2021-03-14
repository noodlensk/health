# frozen_string_literal: true

class CustomFormatter
  attr_reader :io

  def initialize(config)
    @io = config.out_stream

    # Using a block
    config.on_event :test_case_started do |event|
      io.puts "Health check #{event.test_case.name} started"
    end

    config.on_event :test_case_finished do |event|
      log = "Health check #{event.test_case.name} finished"
      log += " with error: '#{event.result.exception}'" if event.result.failed?

      io.puts log
    end
  end
end
