# Simple inventory with basic hosts and groups
vars {
  environment = "test"
  domain      = "example.com"
}

transport "ssh" {
  user = "admin"
  port = 22
  connection_timeout = "30s"
  use_known_hosts = false
}

group "webservers" {
  vars {
    role = "web"
    port = 8080
  }
  host "web1" {
    vars {
      ip = "10.0.1.10"
      hostname = "web1.${vars.domain}"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
  host "web2" {
    vars {
      ip = "10.0.1.11"
      hostname = "web2.${vars.domain}"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
}

group "databases" {
  vars {
    role = "database"
    port = 5432
  }
  host "db1" {
    vars {
      ip = "10.0.2.10"
      hostname = "db1.${vars.domain}"
    }
    transport "ssh" {
      host = "${vars.ip}"
      user = "dbuser"
    }
  }
}
