# Name conflicts between groups and hosts - should fail
group "server1" {
    vars {
        role = "group"
    }
}

host "server2" {
    groups = ["server1"]

    vars {
        ip = "10.0.1.2"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

# This host has same name as group above - should fail
host "server1" {
    vars {
        ip = "10.0.1.1"
        role = "host"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}
