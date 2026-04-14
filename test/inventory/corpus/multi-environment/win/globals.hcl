# Global configuration
vars {
    company = "acme"
    base_domain = "acme.com"
    ssh_user = "deploy"
}

transport "ssh" {
    user = "${var.ssh_user}"
    connection_timeout = "30s"
    use_known_hosts = true
    known_hosts_path = "C:\\temp\\known_hosts"
}
