# Deep parent hierarchy test (5 levels)
vars {
    organization = "acme-corp"
    compliance_level = "high"
}

transport "ssh" {
    user = "admin"
    connection_timeout = "30s"
}

# Level 1: Foundation
group "foundation" {
    vars {
        tier = "foundation"
        security_baseline = "cis-level-2"
        audit_enabled = true
    }
}

# Level 2: Infrastructure
group "infrastructure" {
    parent = "foundation"

    vars {
        tier = "infrastructure"
        monitoring_enabled = true
        backup_retention = "30d"
    }
}

# Level 3: Platform
group "platform" {
    parent = "infrastructure"

    vars {
        tier = "platform"
        container_runtime = "docker"
        orchestrator = "k8s"
    }
}

# Level 4: Application Services
group "app_services" {
    parent = "platform"

    vars {
        tier = "application"
        service_mesh = "istio"
        tracing_enabled = true
    }
}

# Level 5: Specific Application
group "user_service" {
    parent = "app_services"

    vars {
        tier = "microservice"
        service_name = "user-service"
        version = "v2.3.1"
        replicas = 3
    }
}

host "user-svc-1" {
    groups = ["user_service"]

    vars {
        ip = "10.1.1.10"
        instance_id = "i-user-001"
        cpu_cores = 2
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "user-svc-2" {
    groups = ["user_service"]

    vars {
        ip = "10.1.1.11"
        instance_id = "i-user-002"
        cpu_cores = 2
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "user-svc-3" {
    groups = ["user_service"]

    vars {
        ip = "10.1.1.12"
        instance_id = "i-user-003"
        cpu_cores = 4  # Higher capacity instance
    }
    transport "ssh" {
        host = "${var.ip}"
    }
}

# Another Level 5: Different Application
group "payment_service" {
    parent = "app_services"

    vars {
        tier = "microservice"
        service_name = "payment-service"
        version = "v1.8.2"
        replicas = 2
        pci_compliant = true
    }
}

host "payment-svc-1" {
    groups = ["payment_service"]

    vars {
        ip = "10.1.2.10"
        instance_id = "i-payment-001"
        cpu_cores = 4
        secure_enclave = true
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "payment-svc-2" {
    groups = ["payment_service"]

    vars {
        ip = "10.1.2.11"
        instance_id = "i-payment-002"
        cpu_cores = 4
        secure_enclave = true
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

# Separate branch from infrastructure
group "data_platform" {
    parent = "infrastructure"
    vars {
        tier = "data"
        encryption_at_rest = true
        data_classification = "sensitive"
    }
}

group "analytics_db" {
    parent = "data_platform"
    vars {
        tier = "database"
        engine = "postgresql"
        cluster_mode = true
    }
}

host "analytics-db-1" {
    groups = ["analytics_db"]

    vars {
        ip = "10.2.1.10"
        instance_id = "i-analytics-db-001"
        storage_type = "ssd"
        storage_size = "500GB"
        is_primary = true
    }

    transport "ssh" {
        host = "${var.ip}"
        user = "postgres"
    }
}

host "analytics-db-2" {
    groups = ["analytics_db"]

    vars {
        ip = "10.2.1.11"
        instance_id = "i-analytics-db-002"
        storage_type = "ssd"
        storage_size = "500GB"
        is_primary = false
    }

    transport "ssh" {
        host = "${var.ip}"
        user = "postgres"
    }
}

host "analytics-db-3" {
    groups = ["analytics_db"]

    vars {
        ip = "10.2.1.12"
        instance_id = "i-analytics-db-003"
        storage_type = "ssd"
        storage_size = "500GB"
        is_primary = false
    }

    transport "ssh" {
        host = "${var.ip}"
        user = "postgres"
    }
}
