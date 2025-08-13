// Invalid HCL syntax - missing quotes
process {
  name = Invalid Process Name Without Quotes
  targets = "host1"
  
  step "test" {
    name = "Test Step"
    module = "shell"
  }
}
