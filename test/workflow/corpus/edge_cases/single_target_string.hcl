// Single target as string vs array
process {
  name = "Single Target Process"
  targets = "single-host"
  
  step "test" {
    name = "Test Step"
    module = "shell"
    targets = "different-host"
  }
}
