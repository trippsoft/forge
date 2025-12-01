// Step block with too many labels
process {
  name = "Valid Process"
  targets = "host1"
  
  step "id1" "id2" {
    name = "Test Step"
    module = "shell"
  }
}
