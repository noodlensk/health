# frozen_string_literal: true

class CucuGlobals
  include Singleton

  def initialize
    @report = {
      Created: Time.now,
      Took: nil,
      Status: 'FAILED',
      Checks: []
    }
  end

  attr_reader :report
end
