# Parent-only group hierarchy example
vars {
    environment = "production"
    monitoring_enabled = true
}

transport "ssh" {
    user = "deploy"
    connection_timeout = "30s"
    use_known_hosts = true
    known_hosts_path = "/tmp/known_hosts"
}

# Base infrastructure group
group "infrastructure" {
    vars {
        backup_enabled = true
        log_level = "info"
        managed = true
    }

    transport "ssh" {
        user = "infrauser"
    }
}

# Frontend inherits from infrastructure
group "frontend" {
    parent = "infrastructure"

    vars {
        role = "frontend"
        load_balanced = true
        app_port = 8080
    }
}

host "web1" {
    groups = ["frontend"]

    vars {
        ip = "10.0.1.10"
        tier = "primary"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "web2" {
    groups = ["frontend"]

    vars {
        ip = "10.0.1.11"
        tier = "secondary"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "web3" {
    groups = ["frontend"]

    vars {
        ip = "10.0.1.12"
        tier = "secondary"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

# Backend inherits from infrastructure  
group "backend" {
    parent = "infrastructure"

    vars {
        role = "backend"
        api_version = "v2"
        app_port = 9000
    }
}

host "api1" {
    groups = ["backend"]

    vars {
        ip = "10.0.2.10"
        tier = "primary"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "api2" {
    groups = ["backend"]

    vars {
        ip = "10.0.2.11"
        tier = "secondary"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

# Specialized frontend
group "cdn" {
    parent = "frontend"
    vars {
        cache_enabled = true
        edge_locations = ["us-east", "us-west", "eu-west"]
    }
}

host "cdn1" {
    groups = ["cdn"]

    vars {
        ip = "10.0.3.10"
        tier = "edge"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}
