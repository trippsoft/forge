# Circular parent reference - should fail validation
group "group_a" {
    parent = "group_b"

    vars {
        source = "a"
    }
}

host "host1" {
    groups = ["group_a"]

    vars {
        ip = "10.0.1.10"
    }
}

group "group_b" {
    parent = "group_c"

    vars {
        source = "b"
    }
}

group "group_c" {
    parent = "group_a"  # Creates circular reference

    vars {
        source = "c"
    }
}
