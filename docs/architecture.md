# Initial Architecture Sketch

```text
+------------------------+
| Frontend               |
| - job definition       |
| - job submission       |
| - status display       |
+-----------+------------+
            |
            | HTTP / JSON
            v
+-----------+------------+
| API Orchestrator       |
| - request validation   |
| - job creation         |
| - status endpoints     |
+-----+-------------+----+
      |             |
      | SQL         | dispatch jobs
      v             v
+-----+-----+   +---+---------------+
| SQLite DB |   | Compute Worker    |
| - jobs    |   | - execute jobs    |
| - status  |   | - produce results |
+-----------+   +-------------------+

Later:
+------------------------+
| Analysis Service       |
| - convergence analysis |
| - result interpretation|
+------------------------+