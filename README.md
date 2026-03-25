# rbg-api

`rbg-api` is the standalone API and client library for [RoleBasedGroup (RBG)](https://github.com/sgl-project/rbg) — a Kubernetes operator for managing role-based LLM inference workloads (e.g., prefill/decode disaggregation).

This repository is intended to be imported by **upper-layer business systems** that need to interact with RBG resources via the Go client, without depending on the full operator codebase.
