# Transport inheritance example
vars {
  ssh_key_path = "/tmp/test_ssh_key"
  bastion_host = "bastion.example.com"
  admin_user = "admin"
}

# Global SSH transport
transport "ssh" {
  user = "deploy"
  private_key_path = "${vars.ssh_key_path}"
  connection_timeout = "30s"
  use_known_hosts = false
}

# Group with bastion host transport
group "secure_servers" {
  vars {
    security_level = "high"
    compliance_required = true
  }
  transport "ssh" {
    user = "secure_deploy"
    address = "${vars.bastion_host}"
    port = 2222
  }
  host "secure1" {
    vars {
      ip = "10.0.1.10"
      role = "security"
    }    
    # Host-specific transport override
    transport "ssh" {
      host = "${vars.ip}"
      user = "root"
      port = 22
    }
  }
  host "secure2" {
    vars {
      ip = "10.0.1.11"
      role = "security"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
}

# Group inheriting global transport
group "standard_servers" {
  vars {
    security_level = "standard"
    compliance_required = false
  }
  host "web1" {
    vars {
      ip = "10.0.2.10"
      role = "web"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
  host "web2" {
    vars {
      ip = "10.0.2.11"
      role = "web"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
}

# Group with custom transport settings
group "admin_servers" {
  vars {
    admin_access = true
  }
  transport "ssh" {
    user = "${vars.admin_user}"
    port = 22
  }
  host "admin1" {
    vars {
      ip = "10.0.3.10"
      role = "admin"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
}
