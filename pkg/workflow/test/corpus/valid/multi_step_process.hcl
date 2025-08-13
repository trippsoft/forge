// Process with multiple steps and various attributes
process {
  name = "Multi-Step Process"
  targets = ["host1", "host2", "host3"]
  exec_timeout = "5m"
  
  step "setup" {
    name = "Setup Environment"
    module = "shell"
  }
  
  step "configure" {
    name = "Configure System"
    module = "config"
    condition = true
  }
  
  step "deploy" {
    name = "Deploy Application"
    module = "deploy"
    targets = ["host1"]
    exec_timeout = "10m"
  }
}
