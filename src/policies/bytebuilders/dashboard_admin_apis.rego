package policies.bytebuilders.GET.users

default allowed = false
default enable = false
default visible = false

enable {
    input.user.properties.isAdmin == true
}

visible {
    enable
}
