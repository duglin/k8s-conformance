
### [ReplicaSet001](tests/rs.go#L9)

ReplicaSet001 will ...

### [ReplicaSet002](tests/rs.go#L13)

ReplicaSet002 will verify that ...
Conformant implementations MUST ....

### [ReplicaSet003](tests/rs.go#L18)

ReplicaSet003 will verify that ...
Conformant implementations MUST ....

### [ReplicaSet004](tests/rs.go#L23)

ReplicaSet004 will verify that ...
Conformant implementations MUST ....

### [Pod001](tests/pod.go#L11)

Pod001 will verify that simple Pod creation works. The platform MUST
create the specified Pod and queries to retrieve the Pod's metadata MUST
return the same values that were used when it wad created. The Pod
MUST eventually end up in the `Running` state, and then be able to be
deleted. Deleting a Pod MUST remove it from the platform.

Test is run in 'serialize' mode.

### [Pod002](tests/pod.go#L62)

Pod002 will verify that ...
Conformant implementations MUST ....

### [Pod003](tests/pod.go#L67)

Pod003 will do something cool

