## Table of Contents

1 - [Pod](#Pod)
2 - [ReplicaSet](#ReplicaSet)

## Pod

### [Pod001](tests/pod.go#L9)

Pod001 will verify that simple Pod creation works. The platform MUST
create the specified Pod and queries to retrieve the Pod's metadata MUST
return the same values that were used when it wad created. The Pod
MUST eventually end up in the `Running` state, and then be able to be
deleted. Deleting a Pod MUST remove it from the platform


### [Pod002](tests/pod.go#L18)

Pod002 will verify that ...
Conformant implementations MUST ....


### [Pod003](tests/pod.go#L25)



## ReplicaSet

### [ReplicaSet001](tests/rs.go#L9)

ReplicaSet001 will ...


### [ReplicaSet002](tests/rs.go#L15)

ReplicaSet002 will verify that ...
Conformant implementations MUST ....


