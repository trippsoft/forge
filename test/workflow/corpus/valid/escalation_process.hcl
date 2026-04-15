// Process with escalation settings
process {
  name = "Escalated Process"
  targets = "privileged-host"
  
  escalate {
    escalate = true
    impersonate_user = "root"
  }
  
  step "install" {
    name = "Install Package"
    module = "package"
    
    escalate {
      escalate = false
    }
  }
  
  step "start_service" {
    name = "Start Service"
    module = "service"
  }
}
