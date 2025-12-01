// Complex process with all features
process {
  name = "Complex Workflow"
  targets = ["db1", "db2"]
  exec_timeout = "15m"
  
  escalate {
    escalate = true
    impersonate_user = "postgres"
  }
  
  step "backup" {
    name = "Backup Database"
    module = "backup"
    condition = "env.BACKUP_ENABLED == true"
    
    input {
      backup_path = "/backups"
      compress = true
    }
    
    output {
      continue_on_fail = false
      failed_condition = "exitcode != 0"
    }
  }
  
  step "maintenance" {
    name = "Database Maintenance"
    module = "maintenance"
    targets = ["db1"]
    
    loop {
      items = ["vacuum", "reindex", "analyze"]
      label = "operation"
    }
    
    escalate {
      impersonate_user = "dba"
    }
  }
}
