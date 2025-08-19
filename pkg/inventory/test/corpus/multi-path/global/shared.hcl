# Global configuration shared across environments
vars {
    organization = "test-corp"
    domain = "test.local"
    ssh_user = "deploy"
}

transport "ssh" {
    user = "${var.ssh_user}"
    connection_timeout = "30s"
    use_known_hosts = false
}

# Base infrastructure groups
group "infrastructure" {
    vars {
        tier = "infrastructure"
        backup_enabled = true
    }
}
