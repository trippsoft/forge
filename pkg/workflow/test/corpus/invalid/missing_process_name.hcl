// Process missing required name attribute
process {
  targets = "host1"
  
  step "test" {
    name = "Test Step"
    module = "shell"
  }
}
