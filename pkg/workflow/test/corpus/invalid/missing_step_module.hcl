// Step missing required module attribute
process {
  name = "Valid Process"
  targets = "host1"
  
  step "test" {
    name = "Test Step"
  }
}
