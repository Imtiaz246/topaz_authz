package policies.bytebuilders.GET.checkTestOrgType

default isOrgAvailable = false

isOrgAvailable {
    input.user.properties.name = "test_org"
}