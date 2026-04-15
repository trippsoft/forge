// Multiple input blocks (should be error)
process {
  name = "Invalid Process"
  targets = "host1"
  
  step "test" {
    name = "Test Step"
    module = "shell"
    
    input {
      var1 = "value1"
    }
    
    input {
      var2 = "value2"
    }
  }
}
