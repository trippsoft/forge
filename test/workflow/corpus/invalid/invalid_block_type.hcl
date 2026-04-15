// Invalid block type
process {
  name = "Test Process"
  targets = "host1"
  
  invalid_block {
    name = "Invalid Block"
  }
  
  step "test" {
    name = "Test Step"
    module = "shell"
  }
}
