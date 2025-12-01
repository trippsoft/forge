// Step block missing required label
process {
  name = "Valid Process"
  targets = "host1"
  
  step {
    name = "Test Step"
    module = "shell"
  }
}
