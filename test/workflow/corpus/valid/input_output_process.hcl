// Process with input and output blocks
process {
  name = "Input/Output Process"
  targets = "app-server"
  
  step "configure" {
    name = "Configure Application"
    module = "config"
    
    input {
      app_name = "myapp"
      version = "1.2.3"
      config_file = "/etc/myapp/config.yaml"
      debug = true
      port = 8080
    }
    
    output {
      continue_on_fail = false
      changed_condition = "exitcode == 0"
      failed_condition = "exitcode != 0"
    }
  }
}
