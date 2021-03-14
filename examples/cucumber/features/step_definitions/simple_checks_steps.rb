# frozen_string_literal: true

When('I try to execute scenario A') do
  @please_fail = true
end

Then('I should receive a proper response from service') do
  expect(1).to eq(2), 'got wrong response' if @please_fail
rescue RSpec::Expectations::ExpectationNotMetError => e
  @collected_errors.push e.message
end

When('I try to execute scenario B') do ## rubocop:disable Lint/EmptyBlock
end

When('I try to execute scenario C') do  ## rubocop:disable Lint/EmptyBlock
end
