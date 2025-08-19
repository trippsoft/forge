# Variable circular reference - should fail
vars {
    var_a = "${var.var_b}"
    var_b = "${var.var_c}"
    var_c = "${var.var_a}"  # Creates circular dependency
}

host "test" {
    vars {
        ip = "10.0.1.1"
        hostname = "test.${var.var_a}"  # This should fail due to circular reference
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}
