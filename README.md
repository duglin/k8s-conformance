# Kubernetes Conformance Test Suite

This repository contains the Kubernetes conformance test suite.
The purpose of these tests are to ensure that a Kubernetes deployment
behaves properly for a core set of features. This will help ensure that
users can port their applications between Kubernetes deployments and get
consistent results.

### Table of Contents
1. [Conformance Tests](#conformance-tests)
1. [Running the Tests](#running-the-tests)
1. [Structure of the Tests](#structure-of-the-tests)
1. [Contributing](#contributing)
1. [Building the Tool](#building-the-tool)
1. [Reporting Errors in the Tests](#reporting-errors-in-the-tests)
1. [Reporting Errors in a Kubernetes Deployment](#reporting-errors-in-a-kubernetes-deployment)

## Conformance Tests

The conformance test suite is made of a series tests that are executed
against a Kubernetes deployment. These tests are focused on testing from
a user's point of view, therefore they are limited to what can be done
via default Kubernetes command line tooling - such as `kubectl`.

### Running the Tests

The tests are run by executing the `kubecon` program.
This program will assume the following is true:
* `kubectl` is in your PATH
* your environment is setup such that `kubectl` is configured to
  talk to the Kubernetes cluster you want to test. In other words, you
  have your Kubernetes config files and environment variabes setup
  properly so that `kubectl` will work without any additional setup.

To download `kubecon` go to the
[releases](https://github.com/duglin/k8s-conformance/releases) page
and download the appropriate version.

To run `kubecon` simply run it from your command line:
```
$ kubecon
```

The output will be a list of each test that is executed along with the
result - `PASS` or `FAIL`.

TODO: Show Sample Output

### Structure of the Tests

The tests are in the `tests` directory. Typically, tests are grouped
into files such that tests that focus on one particular feature of Kubernetes
will be co-located within the same golang file. But this is not a requirement

Each test MUST have a unique name and SHOULD include some keyword giving
a hint as to the general high-level area of Kubernetes that is being tested.
For example, `Pod001` will focus on testing the management of individual
Pods. The test names are not meant to be very descriptive, rather just provide
a basic hit and a unique number. The numbering does not need to be sequential
and if a test is removed renumbering MUST NOT happen. This will ensure that
people can talk about a particular test by number without fear of it changing
over time.

Each test MUST have a comment block preceeding the test which describes
the purpose of the test and what is expected of the Kubernetes deployment.
This text MUST use [RFC2119](https://www.ietf.org/rfc/rfc2119.txt) language
to be clear what are the mandatory semantics/behavior expected.

For example:
```
// Pod001 will verify that simple Pod creation works. The platform MUST
// create the specified Pod and queries to retrieve the Pod's metadata MUST
// return the same values that were used when it wad created. The Pod
// MUST eventually end up in the `Running` state, and then be able to be
// deleted. Deleting a Pod MUST remove it from the platform.
func Pod001(t *Test) {
    ...
}
```

The comment blocks for each test are extracted and put into a single
[tests.md](tests.md) for easy viewing.

## Contributing

TODO: Add stuff on:
* suggesting new tests
* how to add new tests (PRs)
* how are tests accepted - what's the approval process

## Building the Tool

In order to build the `kubecon` tool you'll need to have `make` and the golang
compiler installed.

* Build `bin/kubecon` for your local platform:
```
$ make
```

* Build `bin/kubecon-OS-ARCH` for all known operating systems/architetures:
```
$ make cross
```

* Clean build environment (erase everything except the href checker tool):
```
$ make clean
```

* Erase all files - including the href checker tool:
```
$ make purge
```

## Reporting Errors in the Tests

Open an [issue](https://github.com/duglin/k8s-conformance/issues).

## Reporting Errors in a Kubernetes Deployment

Talk to the organziation hosting or providing the Kubernetes deployment.

