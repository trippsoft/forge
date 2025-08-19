# Reserved group name - should fail validation
group "all" {  # 'all' is a reserved group name
    vars {
        role = "everything"
    }
}

host "host1" {
    groups = ["all"]

    vars {
        ip = "10.0.1.10"
    }
}
