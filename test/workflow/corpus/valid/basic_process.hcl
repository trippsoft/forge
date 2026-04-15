// Basic process with minimal required attributes
process {
  name = "Basic Process"
  targets = "host1"
  
  step "setup" {
    name = "Setup Environment"
    module = "shell"
  }
}
