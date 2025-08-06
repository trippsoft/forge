# Missing host reference - should fail validation
group "frontend" {
  vars {
    role = "web"
  }
  hosts = ["nonexistent_host"]  # References non-existent host
  host "web1" {
    vars {
      ip = "10.0.1.10"
    }
  }
}
