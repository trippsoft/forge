# Complex variable interpolation examples
vars {
    # Basic variables
    environment = "test"
    domain = "example.com" 
    datacenter = "us-east-1"
    
    # Computed variables
    internal_domain = "internal.${var.domain}"
    external_domain = "external.${var.domain}"

    # Network configuration
    network_prefix = "10.0"
    web_subnet = "${var.network_prefix}.1"
    api_subnet = "${var.network_prefix}.2"
    db_subnet = "${var.network_prefix}.3"
    
    # Application configuration
    app_name = "myapp"
    app_version = "1.2.3"
    app_image = "${var.app_name}:${var.app_version}"

    # Port configuration
    base_port = 8000
    web_port = "${var.base_port + 80}"
    api_port = "${var.base_port + 90}"
    
    # List variables
    availability_zones = ["a", "b", "c"]
    log_levels = ["debug", "info", "warn", "error"]
}

transport "ssh" {
    user = "deploy"
    connection_timeout = "30s"
    use_known_hosts = false
}

group "web" {
    vars {
        role = "web"
        cluster_name = "${var.app_name}-web-${var.environment}"
        service_url = "https://web.${var.external_domain}:${var.web_port}"
        internal_url = "http://web.${var.internal_domain}:${var.web_port}"
    }
}
host "web1" {
    groups = ["web"]

    vars {
        host_id = 1
        ip = "${var.web_subnet}.10"
        hostname = "web${var.host_id}.${var.internal_domain}"
        fqdn = "${var.hostname}"
        zone = "${var.availability_zones[0]}"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "web2" {
    groups = ["web"]

    vars {
        host_id = 2
        ip = "${var.web_subnet}.11"
        hostname = "web${var.host_id}.${var.internal_domain}"
        fqdn = "${var.hostname}"
        zone = "${var.availability_zones[1]}"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

group "database" {
    vars {
        role = "database"
        cluster_name = "${var.app_name}-db-${var.environment}"
        internal_url = "postgres://db.${var.internal_domain}:5432"
    }
}

host "db1" {
    groups = ["database"]

    vars {
        host_id = 1
        ip = "${var.db_subnet}.10"
        hostname = "db${var.host_id}.${var.internal_domain}"
        fqdn = "${var.hostname}"
        zone = "${var.availability_zones[0]}"
        role = "primary"
    }

    transport "ssh" {
        host = "${var.ip}"
        user = "postgres"
    }
}
