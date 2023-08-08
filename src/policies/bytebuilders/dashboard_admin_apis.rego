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

#isBackendDeveloper {
#    ds.relation({
#        "object": {
#            "key": "appscode-backend-team",
#            "type": "group"
#        },
#        "realtion": {
#            "name": "member",
#            "object_type": "group"
#        },
#        "subject": {
#            "key": input.user.key,
#            "type": input.user.type
#        }
#    })
#}