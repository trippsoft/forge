// Empty escalate block
process {
  name = "Empty Escalate Process"
  targets = "host1"
  
  escalate {
  }
  
  step "test" {
    name = "Test Step"
    module = "shell"
  }
}
