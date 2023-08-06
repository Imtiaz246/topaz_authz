package b3.TEST.rego

check_relation {
    ds.relation({
      "object": {
          "type": "group",
          "key": "test_org"
        },
      "relation": {
        "name": "member",
        "object_type": "group"
      },
      "subject": {
        "key": "imtiaz@appscode.com",
        "type": "user"
      }
    })
}