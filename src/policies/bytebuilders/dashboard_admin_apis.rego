package policies.bytebuilders.GET.users

import data.ds.relation

default allowed = false
default enable = false
default visible = false

enable {
    input.user.properties.isAdmin == true
}

visible {
    enable
}

isBackendDeveloper {
    ds.relation({
        "object": {
            "key": input.user.key,
            "type": input.user.type
        },
        "relation": {
            "name": "editor",
            "object_type": "bytebuilders.organization"
        },
        "subject": {
            "key": "test_org",
            "type": "bytebuilders.organization"
        }
    })
}