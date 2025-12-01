// Loop block missing required items attribute
process {
  name = "Invalid Process"
  targets = "host1"
  
  step "test" {
    name = "Test Step"
    module = "shell"
    
    loop {
      label = "item"
    }
  }
}
