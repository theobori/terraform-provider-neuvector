resource "neuvector_admission_rule" "containers_prevention" {
    rule_id = 666
    rule_type = "deny"
    category = "Kubernetes"
    comment = "Containers prevention"
    
    criteria {
        name = "runAsRoot"
        op = "="
        path = "runAsRoot"
        value = "true"
    }

    criteria {
        name = "runAsPrivileged"
        op = "="
        path = "runAsPrivileged"
        value = "true"
    }

    disable = false
}
