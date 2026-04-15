// Multiple output blocks (should be error)
process {
  name = "Invalid Process"
  targets = "host1"
  
  step "test" {
    name = "Test Step"
    module = "shell"
    
    output {
      continue_on_fail = false
    }
    
    output {
      changed_condition = "true"
    }
  }
}
