// Multiple escalate blocks (should be error)
process {
  name = "Invalid Process"
  targets = "host1"
  
  escalate {
    escalate = true
  }
  
  escalate {
    impersonate_user = "root"
  }
  
  step "test" {
    name = "Test Step"
    module = "shell"
  }
}
