# Complex variable interpolation examples
vars {
  # Basic variables
  environment = "test"
  domain = "example.com" 
  datacenter = "us-east-1"
  
  # Computed variables
  internal_domain = "internal.${vars.domain}"
  external_domain = "external.${vars.domain}"

  # Network configuration
  network_prefix = "10.0"
  web_subnet = "${vars.network_prefix}.1"
  api_subnet = "${vars.network_prefix}.2"
  db_subnet = "${vars.network_prefix}.3"
  
  # Application configuration
  app_name = "myapp"
  app_version = "1.2.3"
  app_image = "${vars.app_name}:${vars.app_version}"

  # Port configuration
  base_port = 8000
  web_port = "${vars.base_port + 80}"
  api_port = "${vars.base_port + 90}"
  
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
    cluster_name = "${vars.app_name}-web-${vars.environment}"
    service_url = "https://web.${vars.external_domain}:${vars.web_port}"
    internal_url = "http://web.${vars.internal_domain}:${vars.web_port}"
  }
  host "web1" {
    vars {
      host_id = 1
      ip = "${vars.web_subnet}.10"
      hostname = "web${vars.host_id}.${vars.internal_domain}"
      fqdn = "${vars.hostname}"
      zone = "${vars.availability_zones[0]}"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
  host "web2" {
    vars {
      host_id = 2
      ip = "${vars.web_subnet}.11"
      hostname = "web${vars.host_id}.${vars.internal_domain}"
      fqdn = "${vars.hostname}"
      zone = "${vars.availability_zones[1]}"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
}

group "database" {
  vars {
    role = "database"
    cluster_name = "${vars.app_name}-db-${vars.environment}"
    internal_url = "postgres://db.${vars.internal_domain}:5432"
  }
  host "db1" {
    vars {
      host_id = 1
      ip = "${vars.db_subnet}.10"
      hostname = "db${vars.host_id}.${vars.internal_domain}"
      fqdn = "${vars.hostname}"
      zone = "${vars.availability_zones[0]}"
      role = "primary"
    }
    transport "ssh" {
      host = "${vars.ip}"
      user = "postgres"
    }
  }
}
