# Simple Config Example

This directory contains the minimal config needed to use DPT. It defines a
simple project that will store its audit logs in the bigquery dataset
`example_project_logs` and its Terraform state in the storage bucket
`example_project_state` in the same project.

This is a good example for those new to using DPT and can be used to create test
project. You can incrementally add resources to test how they are deployed. In a
production project, you will likely want to setup central auditing, devops and
monitoring. For a complete example, see the `full` sample.
