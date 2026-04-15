// Multiple loop blocks (should be error)
process {
  name = "Invalid Process"
  targets = "host1"
  
  step "test" {
    name = "Test Step"
    module = "shell"
    
    loop {
      items = ["a", "b"]
      label = "item"
    }
    
    loop {
      items = ["x", "y"]
      label = "other"
    }
  }
}
