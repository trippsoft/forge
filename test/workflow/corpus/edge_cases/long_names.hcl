// Very long names and identifiers
process {
  name = "A Very Long Process Name That Tests The Parser's Ability To Handle Extended Names Without Issues"
  targets = "host-with-very-long-name-that-should-still-be-parsed-correctly"
  
  step "step_with_very_long_identifier_name_that_tests_parser_limits" {
    name = "A Very Long Step Name That Also Tests The Parser's String Handling Capabilities"
    module = "module-with-long-name"
  }
}
