# Complex variable interpolation and inheritance
vars {
    global_env = "staging"
    base_domain = "internal.company.com"
    backup_schedule = "0 2 * * *"
}

transport "ssh" {
    user = "deploy"
    connection_timeout = "30s"
}

# Base configuration group
group "base" {
    vars {
        monitoring = true
        log_retention_days = 30
        environment = "${var.global_env}"
        fqdn_suffix = "${var.base_domain}"
    }
}

# Application tier inheriting from base
group "app_tier" {
    parent = "base"

    vars {
        app_version = "2.1.0"
        health_check_port = 8080
        service_url = "https://app.${var.fqdn_suffix}"
    }
}

# Frontend services
group "frontend" {
    parent = "app_tier"

    vars {
        role = "frontend"
        load_balancer_pool = "frontend-${var.environment}"
        replicas = 3
    }
}

host "web1" {
    groups = ["frontend"]

    vars {
        ip = "10.0.1.10"
        hostname = "web1.${var.fqdn_suffix}"
        server_id = 1
        memory_limit = "2GB"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "web2" {
    groups = ["frontend"]

    vars {
        ip = "10.0.1.11"
        hostname = "web2.${var.fqdn_suffix}"
        server_id = 2
        memory_limit = "2GB"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "web3" {
    groups = ["frontend"]

    vars {
        ip = "10.0.1.12"
        hostname = "web3.${var.fqdn_suffix}"
        server_id = 3
        memory_limit = "4GB"  # More memory for this instance
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

# Backend services  
group "backend" {
    parent = "app_tier"

    vars {
        role = "backend"
        database_url = "postgres://db.${var.fqdn_suffix}:5432/app"
        api_prefix = "/api/v2"
    }
}

host "api1" {
    groups = ["backend"]

    vars {
        ip = "10.0.2.10"
        hostname = "api1.${var.fqdn_suffix}"
        instance_type = "primary"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "api2" {
    groups = ["backend"]

    vars {
        ip = "10.0.2.11"
        hostname = "api2.${var.fqdn_suffix}"
        instance_type = "secondary"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

# Data tier with different inheritance
group "data_tier" {
    parent = "base"

    vars {
        role = "database"
        backup_enabled = true
        replication_factor = 2
    }
}

group "databases" {
    parent = "data_tier"

    vars {
        engine = "postgresql"
        version = "14.9"
        port = 5432
    }
}

host "db1" {
    groups = ["databases"]

    vars {
        ip = "10.0.3.10"
        hostname = "db1.${var.fqdn_suffix}"
        is_primary = true
        storage_size = "100GB"
    }

    transport "ssh" {
        host = "${var.ip}"
        user = "postgres"
    }
}

host "db2" {
    groups = ["databases"]

    vars {
        ip = "10.0.3.11"
        hostname = "db2.${var.fqdn_suffix}"
        is_primary = false
        storage_size = "100GB"
    }

    transport "ssh" {
        host = "${var.ip}"
        user = "postgres"
    }
}
