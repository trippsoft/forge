# Invalid parent reference - should fail validation
group "frontend" {
    parent = "nonexistent_group"  # References non-existent group

    vars {
        role = "web"
    }
}

host "web1" {
    groups = ["frontend"]

    vars {
        ip = "10.0.1.10"
    }
}
