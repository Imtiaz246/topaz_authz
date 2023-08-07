# topaz_authz

1. To configure `resource` and `name` using topaz cli: 
    - `topaz configure -d -s -r ghcr.io/aserto-policies/policy-todo-rebac:latest -n policy-todo`
2. To import some static pre-defined user data
    - Make a file json file with `objects` array and another one with `relations` array
    - `topaz import -i -d .` `to import from current directory`
3. API call to get the policies: 
    - `curl -k https://localhost:8383/api/v2/policies`

# Authorizer
1. The authorizer uses the OPA to compute a decision based on a `policy`, `user-context` and `data`
2. The authorizer exposes 3 api endpoints to calling application
   - For Authorization - authz
   - For getting all policies applied to the system - policies
   - For getting build information - info
   
   And all are exposed with `/api/v2` prefix

# Authorization
1. The authorizer exposes 3 api endpoints available for the calling application
    - `authz/is`
    - `authz/query`
    - `ahthz/decisiontree`
   
   All of these are `POST` style API and accepts a body or request payload
   
   The payload for each of these APIs includes 3 common objects:
   - `Identity context` - identifies the user
   - `Policy context` - identifies the policy
   - `Resource context` - identifies the resource (optional)

# Identity context
1. Identity context is to identify the user. The `type` property of that identity context `tells the authorizar` how to authorize the user
   - `IDENTITY_TYPE_NONE`: the user context is empty (`input.user` is empty) 
   - `IDENTITY_TYPE_SUB`: the user context is passed as an OAUTH subject
   - `IDENTITY_TYPE_JWT`: the user context is passed as a JWT access token. In this case the authorizer automatically will extract the `sub`(subject) claim from the token.
   
   The authorizer will extract the `subject` regarding the identity type and use the subject to look up for the user in its local directory. If found it'll load the user object identified by the identity context
   and make it available to the policy as `input.user`
   ```json
    {
      "identityContext": {
        "identity": ["imtiazCho"],
        "type": "IDENTITY_TYPE_SUB"
      }    
    }    
   ```
   
# Resource context
1. The `resoruceContext` is a key-value map that is passed into the authorizer and materialized as `input.resource` in the policy.
   ```json
   {
      "resourceContext": {
          "ownerKey": ["owner-key"]
      }
   }
   ```
   
# Policy context
1. The policy context identifies the policy and decision(s) to evaluate. It takes `path` for package bundle and `decisions` array which denotes one or more decisions to be made by the authorizer. The `policyContext` passed in will be available to the policy as `input.policy`.
   ```json
   {
      "policyContext": {
          "decisions": [
              "allowed"
          ],
          "path": "bb.GET.api.admin"
      }  
   }
   ```
   The policy context above will evaluate the `allowed` decision for the policy module `bb.GET.api.admin`
2. The common usage for `policyContext` in the `decisiontree` API is to identify the policy ID and policy root
   - `POST .../api/v2/authz/decisiontree`
   ```json
   {
       "policyContext": {
           "decisions": [
               "visible",
               "enabled"
           ],
           "path": "sample"
       } 
   }
   ```
   This call will evaluate all paths under the "sample" root, and return the values of the "visible" and "enabled" decisions using the `identityContext` and `resourceContext` that may also be passed in.

# is API
1. The `is` API is the primary API for determining whether ia user is authorized to perform an action on a resource.
   ```json
   {
      "identity_context": {
          "type": "IDENTITY_TYPE_SUB",
          "identity": "rick@the-citadel.com"
      },
      "policy_context": {
           "path": "todoApp.GET.todos",
           "decisions": ["allowed"]
       },
       "resourceContext": {
           "additionalProp1": "string",
           "additionalProp2": "string",
           "additionalProp3": "string"
       }
   }
   ```
   
# authz/decisiontree
1. The `decissiontree` API allows the caller to get the value of any decisions across ALL policy modules, with a user context , but without a resource context.
   The API is useful for getting a `decission tree` that guides a calling application around what functionality will be available to a user based on their context.
   ```json
   {
       "identityContext": {
           "type": "IDENTITY_TYPE_SUB",
           "identity": "<subject>"
       },
       "policyContext": {
           "decisions": ["visible", "enabled"],
           "path": "sample"
       },
       "resourceContext": {
           "additionalProp1": "string",
           "additionalProp2": "string",
           "additionalProp3": "string"
       },
       "options": {
           "pathSeparator": "PATH_SEPARATOR_SLASH"
       }
   }
   ```
   The `options` map allows the caller to specify the format for retrieving the cartesian product of paths and decisions that are being requested.

# Directory
1. Topaz directory stores 3 types of entity: 
  - Objects: Represents the participants in the authorization decision. Some subjects are subjects and some are resources.
  - Permissions: an action that `subjects (user)` may attempt to perform on objects
  - Relation: a labeled association between a source object(resource) and a target object(subject)
2. To help define different kinds of objects and specify the relations between them, topaz directory provides 2 extensible type of types: 
   - Objects type: defines the kinds of objects (including subjects) that can be created in the directory
     - Built in objects types: (User, Group, Identity, Application, Resource)
   - Relations: defines the relations that can be created between directory objects
3. An object is identified by the combination of `object-type` and `object-key`
4. A relation type is uniquely identified by the combination of `object-type-name` and `relation-name`
        