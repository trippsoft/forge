// Empty input block
process {
  name = "Empty Input Process"
  targets = "host1"
  
  step "test" {
    name = "Test Step"
    module = "shell"
    
    input {
    }
  }
}
