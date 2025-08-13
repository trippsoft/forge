// Invalid HCL syntax - missing closing brace
process {
  name = "Invalid Syntax Process"
  targets = "host1"
  
  step "test" {
    name = "Test Step"
    module = "shell"
  // Missing closing brace
