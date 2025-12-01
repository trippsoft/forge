// Process with loop functionality
process {
  name = "Looping Process"
  targets = ["web1", "web2", "web3"]
  
  step "deploy_services" {
    name = "Deploy Services"
    module = "deploy"
    
    loop {
      items = ["nginx", "redis", "postgres"]
      label = "service"
      condition = true
    }
  }
}
