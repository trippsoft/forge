// Process block with invalid labels
process "invalid_label" {
  name = "Invalid Process"
  targets = "host1"
  
  step "test" {
    name = "Test Step"
    module = "shell"
  }
}
