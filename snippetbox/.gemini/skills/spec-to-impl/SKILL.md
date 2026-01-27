---
name: spec-to-impl
description: Convert high-level requirements into technical implementation plans.
---
# Specification-to-Implementation Protocol

When a new feature or requirement is described:
1. **Plan First**: Do not write code immediately. Provide a high-level plan listing the files to be modified and the logic flow.
2. **Review Diffs**: Use clear Diff views to show exactly what will be added, modified, or deleted.
3. **Side-Effect Analysis**: Proactively warn the user if the change affects Database Schemas, environment variables, or existing API contracts.
4. **Scalability Check**: Briefly mention if the proposed implementation might have performance bottlenecks (e.g., N+1 queries).