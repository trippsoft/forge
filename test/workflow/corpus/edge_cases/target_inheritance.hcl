// Target inheritance from process to steps
process {
  name = "Target Inheritance Process"
  targets = ["parent1", "parent2"]
  exec_timeout = "30s"
  
  step "inherits_targets" {
    name = "Step That Inherits Targets"
    module = "shell"
    // No targets specified - should inherit from process
  }
  
  step "overrides_targets" {
    name = "Step That Overrides Targets"
    module = "shell"
    targets = ["override1"]
  }
  
  step "inherits_timeout" {
    name = "Step That Inherits Timeout"
    module = "shell"
    // No exec_timeout specified - should inherit from process
  }
}
