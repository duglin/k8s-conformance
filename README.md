# Kubernetes Conformance Test Suite

## Abstract

This repository contains the Kubernetes conformance test suite.
The purpose of these tests are to ensure that a Kubernetes deployment
behaves properly for a core set of features. This will help ensure that
users can port their applications between Kubernetes deployments and get
consistent results.

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

To download `kubedown` go to: ...

To run `kubecon` simply run it from your command line:
```
$ kubecon
```

The output will be a list of each test that is executed along with the
result - `PASS` or `FAIL`.

*** Show Sample Output ***

### Structure of the Tests

The tests are in the `tests` directory. Typically, tests are grouped
into files such that tests that focus on one particular feature of Kubernetes
will be co-located within the same golang file. But this is not a requirement

Each test MUST have a unique name and SHOULD include some keyword giving
a hint as to the general high-level area of Kubernetes that is being tested.
For example, `Pod001` will focus on testing the management of individual
Pods. The test names are not means to be very descriptive, rather just provide
a basic hit and a unique number. The numbering does not need to be sequential
and if a test is removed renumbering MUST NOT happen. This will ensure that
people can talk about a particular test by number without fear of it changing
over time.

Each test MUST have a comment block preceeding the test which describes
the purpose of the test and what is expected of the Kubernetes deployment.
This text MUST use [RFC2119](https://www.ietf.org/rfc/rfc2119.txt) language
to be clear what are the mandatory semantics/behavior expected.

The comment blocks for each test are extracted and put into a single
[tests.md](tests.md) for easy viewing.

## Contributing

Add stuff on:
* suggesting new tests
* how to add new tests (PRs)
* how are tests accepted - what's the approval process

## Reporting Errors in a Kubernetes Deployment

Talk to the organziation hosting or providing the Kubernetes deployment.

## Reporting Errors in the Tests

Open an [issue](issues).

