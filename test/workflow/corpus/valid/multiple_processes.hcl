// Multiple processes in one file
process {
  name = "First Process"
  targets = "host1"
  
  step "step1" {
    name = "First Step"
    module = "shell"
  }
}

process {
  name = "Second Process"  
  targets = ["host2", "host3"]
  
  step "step1" {
    name = "Another Step"
    module = "copy"
  }
  
  step "step2" {
    name = "Final Step"
    module = "service"
  }
}
