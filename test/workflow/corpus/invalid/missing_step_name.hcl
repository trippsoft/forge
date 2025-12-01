// Step missing required name attribute
process {
  name = "Valid Process"
  targets = "host1"
  
  step "test" {
    module = "shell"
  }
}
