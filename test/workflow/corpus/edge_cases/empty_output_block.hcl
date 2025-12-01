// Empty output block
process {
  name = "Empty Output Process"
  targets = "host1"
  
  step "test" {
    name = "Test Step"
    module = "shell"
    
    output {
    }
  }
}
