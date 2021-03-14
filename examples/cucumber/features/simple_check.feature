Feature: Simple health check

  In order to prevent downtime
  As a software engineer
  I want be able to run health checks against my service

  @only-a
  Scenario: Simple check scenario A
    When I try to execute scenario A
    Then I should receive a proper response from service

  Scenario: Simple check scenario B
    When I try to execute scenario B
    Then I should receive a proper response from service

  Scenario: Simple check scenario C
    When I try to execute scenario C
    Then I should receive a proper response from service