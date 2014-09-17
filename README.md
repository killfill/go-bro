
# Cloud Foundry SQL Broker

This is an experimental SQL broker for Cloud foundry, written in GO.

It implements the Broker API v2, and only works with PostgreSQL.


## Configurations

The file *catalog.json* contains the catalog offerings

The file *config.json* contains the broker configuration, auth and plans.
The plans of config and catalog should match by its uuid...


## TODO

* Auth requests
* Support plan limits
* MySQL support

