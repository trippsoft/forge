// Process missing required targets attribute
process {
  name = "Missing Targets Process"
  
  step "test" {
    name = "Test Step"
    module = "shell"
  }
}
