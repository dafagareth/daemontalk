---
title: "DuckDB In-Process Vectorized Execution on Parquet and Arrow"
slug: d8e7f6a5
aliases: [duckdb-analytical-queries]
date: 2026-06-25
tags: [storage, tools, backend]
lang: en
draft: false
type: post
---

DuckDB is an embeddable analytical SQL database engine optimized for columnar processing workloads. Traditional row-oriented OLTP databases struggle with large aggregations because tuple overhead forces irrelevant attributes into CPU cache. DuckDB uses a vectorized engine that operates on contiguous memory blocks, executing fast queries over local and remote Parquet or CSV datasets.

## Fun Facts

**Fact 1.** DuckDB was created by Hannes Mühleisen and Mark Raasveldt at the Centrum Wiskunde & Informatica in 2019 to provide an analytical counterpart to SQLite.

**Fact 2.** DuckDB uses morsel-driven parallelism to divide analytical tasks into small work units assigned dynamically to available CPU threads without client-server IPC overhead.

**Fact 3.** DuckDB queries remote S3 Parquet files by issuing HTTP range requests, downloading only the footer metadata and requested column chunks rather than whole files.

---

## Tips and Tricks

### 1. Query Parquet and CSV Files Directly

You can execute analytical SQL queries over raw Parquet and compressed CSV files without pre-loading them into database tables.

```sql
-- Read multiple Parquet files with wildcard path matching
SELECT 
    category,
    COUNT(*) AS total_orders,
    ROUND(AVG(amount), 2) AS avg_amount
FROM read_parquet('data/orders_2026_*.parquet')
WHERE status = 'COMPLETED'
GROUP BY category
ORDER BY total_orders DESC;
```

### 2. Zero-Copy Data Exchange with Apache Arrow in Python

DuckDB integrates with Apache Arrow in Python, allowing memory buffers to be shared without serialization or data copying.

```python
import duckdb
import pyarrow as pa

# Create an in-memory PyArrow table
data = pa.table({'id': [1, 2, 3], 'val': [10.5, 20.0, 30.2]})

# Query Arrow table directly using DuckDB zero-copy conversion
con = duckdb.connect()
result = con.execute("SELECT AVG(val) FROM data WHERE id > 1").fetchone()
print(f"Average value: {result[0]}")
```

### 3. Run In-Process Analytical Queries in Go

You can integrate DuckDB into Go applications using the `database/sql` driver for low-latency analytical data processing.

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/marcboeker/go-duckdb"
)

func main() {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var sum float64
	err = db.QueryRow("SELECT SUM(amount) FROM read_csv_auto('sales.csv')").Scan(&sum)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total sales sum: %.2f\n", sum)
}
```

### 4. Tune Parallel Threads and Memory Limits

Control memory allocation and thread concurrency to optimize resource consumption during large aggregate joins.

```sql
-- Limit memory usage to 4 gigabytes and set worker threads
SET max_memory = '4GB';
SET threads = 4;

-- Verify active memory settings
SELECT name, value FROM duckdb_settings() WHERE name LIKE '%memory%';
```

### 5. Export Query Results to Compressed Parquet Files

Persist aggregated analytical results directly into Snappy-compressed Parquet files for distribution.

```sql
-- Export query output directly to Parquet format
COPY (
    SELECT user_id, COUNT(*) AS login_count
    FROM 'logs/*.parquet'
    GROUP BY user_id
) TO 'output/user_summary.parquet' (FORMAT PARQUET, COMPRESSION SNAPPY);
```
