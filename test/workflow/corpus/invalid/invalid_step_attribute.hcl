// Invalid attribute in step
process {
  name = "Test Process"
  targets = "host1"
  
  step "test" {
    name = "Test Step"
    module = "shell"
    invalid_attribute = "should not be here"
  }
}
