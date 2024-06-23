---
title: Migrating from SQLite to CockroachDB
description: Small article that review the migration from SQLite to CockroachDB.
tags: ['database', 'sqlite', 'cockroachdb', 'devops']
---

## Table of contents

<div class="toc">

\\{\\{ $.TOC }}

</div>

## Introduction

The last 3 days, I have been working on migrating my application from SQLite to CockroachDB. The reason is that I'm actually using NFS to share the SQLite database between my servers, and it's not really efficient. Grafana started to show serious "sync" issues, meaning that the database was locked for too long, causing my alerts to go wild.

This article is a simple review of the migration to CockroachDB.

## What is CockroachDB?

CockroachDB is a **distributed** SQL database built on a transactional and strongly-consistent key-value store. The syntax is similar to PostgreSQL, but some features are missing, and some defaults are different (It's saner in my opinion).

CockroachDB offers quite a lot of features, such as:

- **Distributed SQL**: CockroachDB is a distributed SQL database that scales horizontally using the Raft consensus algorithm to ensure that data is always consistent.
- **Geo-partitioning**: CockroachDB allows you to partition your data based on the location of your users, ensuring that data is always close to where it's needed.
- **Automatic failover**: CockroachDB automatically handles failover, ensuring that your data is always available.
- **Backup and restore**: CockroachDB offers a simple backup and restore mechanism that allows you to easily recover from data loss over S3.

## Missing features and critical differences compared to PostgreSQL

To migrate from SQLite to CockroachDB, you should know there are some critical features missing compared to PostgreSQL:

- **`ALTER TABLE ... ALTER COLUMN ... TYPE`**: This feature is experimental and can be enabled with `enable_experimental_alter_column_type_general=true` as a session variable. This feature **does not work with transactions**. You may need to manually update your schema.
- The default integer type is `INT8` instead of `INT4` in PostgreSQL, which can cause conversion issues (like with Rust or JS). To set the default behavior, use `default_int_size=4` as a session variable.
- **`SERIAL`** is not same as in PostgreSQL. CockroachDB defaults to `INT8` and doesn't use sequences. Instead, it calls `unique_rowid()` to generate a unique ID. To set the default behavior, use `serial_normalization=sql_sequence_cached` as session variable.
- Locale "C" is not supported.
- A `rowid` column is automatically created if the table doesn't have a primary key. This can conflict with some migration steps (like adding a primary key later).

## Deployment

CockroachDB deployment is quite simple. Today, CockroachDB recommends using the [`cockroach-operator`](https://github.com/cockroachdb/cockroach-operator) to deploy CockroachDB on Kubernetes. The operator is a Kubernetes operator that manages CockroachDB clusters. It automates tasks like scaling, backups, and upgrades.

While it's quite powerful, I prefer using the Helm chart, which is more explicit and easier to manage. The Helm chart is available on the [official Helm repository](https://github.com/cockroachdb/helm-charts).

Since the deployment is quite simple, I won't go into details. You can find the official documentation on the [CockroachDB website](https://www.cockroachlabs.com/docs/stable/kubernetes-overview). Just remember to use `cert-manager` instead of their signer.

## Migration

The migration steps are:

1. Dump the SQLite database:

   ```shell
   sqlite3 mydb.sqlite

   .output mydb.sql
   .dump
   .quit
   ```

2. Create a new database and its user in Cockroach:

   ```shell
   cockroach sql --certs-dir=certs --host=my-cockroachdb-public
   ```

   ```sql
   CREATE DATABASE mydb;
   CREATE USER myuser WITH PASSWORD 'mypassword';
   GRANT ALL ON DATABASE mydb TO myuser;
   ```

3. Apply compatibility flags:

   ```sql
   ALTER ROLE grafana SET enable_experimental_alter_column_type_general=true;
   ALTER ROLE grafana SET default_int_size=4;
   ALTER ROLE grafana SET serial_normalization=sql_sequence_cached;
   ```

4. Start the application to be migrated with PostgreSQL:

   ```shell
   export DATABASE_URL=postgresql://myuser:mypassword@my-cockroachdb-public:26257/mydb?sslmode=require
   ```

5. Wait for crashes and fix them.

6. Import the data from the SQL dump. Basically, edit the file `mydb.sql` to only have the `INSERT` statements. Check also for `BOOL` since SQLite doesn't support it and use integers instead. Replace the value with `TRUE` or `FALSE`. After that, runs the SQL commands manually to insert the data.

## Adding monitoring

CockroachDB offers monitoring out of the box by using Prometheus and the `ServiceMonitor` CRD. You can find the official documentation on the [CockroachDB website](https://www.cockroachlabs.com/docs/stable/monitor-cockroachdb-with-prometheus).

## Adding backup

Since we are now using a complex database, it is good practice to back up regularly. CockroachDB offers a simple backup and restore mechanism that allows you to easily recover from data loss over S3.

I'm using a S3-compatible storage, therefore I can use the following command to backup my database:

```sql
CREATE EXTERNAL CONNECTION backup_bucket AS 's3://my-bucket?AWS_ACCESS_KEY_ID=my-access-key&AWS_SECRET_ACCESS_KEY=my-secret-key&AWS_REGION=my-region&AWS_ENDPOINT=my-endpoint';
BACKUP INTO 'external://backup_bucket';
```

To schedule backups, CockroachDB is able to schedule:

```sql
CREATE SCHEDULE backup_schedule FOR BACKUP INTO 'external://backup_bucket' RECURRING '@daily' WITH SCHEDULE OPTIONS first_run='now';
```

## Last improvements

On the S3 provider, remember to add lifecycle rules to archive the backups and change the storage class to `GLACIER` after a certain period. This will save you a lot of money.

## What has been migrated? What couldn't be migrated?

Grafana was the first application to be migrated. The issues are already been cited:

- **`ALTER TABLE ... ALTER COLUMN ... TYPE`** not working with transactions.
- **`SERIAL`** not the same as in PostgreSQL.
- Integer type default is `INT8` instead of `INT4`.
- Locale "C" is not supported.

There is still one issue about `ctid` not being found.

Second application to be migrated: VaultWarden. The issues are the same as Grafana.

Lastly, the blog, which was developed with CockroachDB in mind.

One thing that couldn't be migrated is Joplin, the note-taking application. Let me lash out on this product for a moment:

- Coded in TS and it LAGS like hell.
- Hard-coded PostgreSQL configuration compared to Knex.js.
- SQLite database is 5000MB?! While there is only 12MB of data?!
- After rebooting back to SQLite ('cause I couldn't migrate to CockroachDB), the password is denied?!

Nope, f- Joplin. I'm going with Obsidian + Syncthing.

## Conclusion

CockroachDB is cool, cool like cool for kids. To deploy CockroachDB, I didn't need to configure much, but the migration was a bit painful, and I'm not sure if there are any runtime issues.

Monitoring-wise:

![image-20240623162810084](./page.assets/image-20240623162810084.png)

The database seems to reach around 2GB of memory, with Grafana being the intensive player. You can see that the `cockroach-client-secure` is leaking, I had to kill it.

Other than that, I can just the that the dashboards are pretty cool:

![image-20240623163235631](./page.assets/image-20240623163235631.png)

![image-20240623163259218](./page.assets/image-20240623163259218.png)

![image-20240623163322218](./page.assets/image-20240623163322218.png)

Overall, yeah, it's cool. Grafana is stable, VaultWarden is replicated, CockroachDB is replicated and backed up. I guess I reached a good level of stability.
