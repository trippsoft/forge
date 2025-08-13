// Comments and whitespace handling
process {
  // Process comment
  name = "Comment Test Process" // inline comment
  targets = "host1"
  /* Multi-line
     comment */
  
  step "test" {
    name = "Test Step"
    module = "shell"
    
    // Step comment
    condition = true /* another inline comment */
  }
  
  // Another step
  step "another" {
    name = "Another Step"
    module = "copy"
  }
}
