---
subcategory: "Distributed Cache Service (DCS)"
---

# sbercloud_dcs_parameters

Manages a DCS configuration parameters within SberCloud.

## Example Usage

```hcl
variable instance_id {}
variable project_id {}

resource "sbercloud_dcs_parameters" "config_1" {
  instance_id = var.instance_id
  project_id  = var.project_id

  parameters = {
    timeout = "1000"
    maxclients = "2100"
    appendfsync = "always"
    maxmemory-policy = "allkeys-random"
    zset-max-ziplist-value = "128"
    repl-timeout = "120"
  }
}
```

## Argument Reference

The following arguments are supported:

*  `instance_id` - (Required, String) Specifies the ID of the instance.
*  `project_id` - (Required, String) Specifies the project.
*  `parameters` - (Required, Map) A mapping of parameters to assign to the DCS instance. 
   Each parameter is represented by one key-value pair.
   + `timeout` - (Optional, String) Close the connection after a client is idle for N seconds (0 to disable). 
   Value range: 0-7200. Default value: 0. **Works on Redis & Memcached.**
   + `maxmemory-policy` - (Optional, String) How Redis will select what to remove when maxmemory is reached,
   You can select among five behaviors: volatile-lru : remove the key with an expire set using an LRU algorithm;
   allkeys-lru: remove any key according to the LRU algorithm; volatile-lfu:remove the key with an expire set using an LFU algorithm;
   allkeys-lfu:remove any key according to the LFU algorithm; volatile-random: remove a random key with an expire set;
   allkeys-random: remove a random key, any key; volatile-ttl: remove the key with the nearest expire time (minor TTL);
   noeviction: don't expire at all, just return an error on write operations. Default value: volatile-lru. **Works on Redis & Memcached.**
   + `hash-max-ziplist-entries` - (Optional, String) Hashes are encoded using a memory efficient data structure
   when they have a small number of entries. Value range: 1-10000. Default value: 512. **Works only on Redis.**
   + `hash-max-ziplist-value` - (Optional, String) Hashes are encoded using a memory efficient data structure 
   when the biggest entry does not exceed a given threshold. Value range: 1-10000. Default value: 64. **Works only on Redis.**
   + `set-max-intset-entries` - (Optional, String) When a set is composed of just strings that happen to be integers
   in radix 10 in the range of 64 bit signed integers. Value range: 1-10000. Default value: 512. **Works only on Redis.**
   + `zset-max-ziplist-entries` - (Optional, String) Sorted sets are encoded using a memory efficient data structure
   when they have a small number of entries. Value range: 1-10000. Default value: 128. **Works only on Redis.**
   + `zset-max-ziplist-value` - (Optional, String) Sorted sets are encoded using a memory efficient data structure when
   the biggest entry does not exceed a given threshold. Value range: 1-10000. Default value: 64. **Works only on Redis.**
   + `latency-monitor-threshold` - (Optional, String) Only events that run in more time than the configured latency-monitor-threshold
   will be logged as latency spikes. If latency-monitor-threshold is set to 0, latency monitoring is disabled. 
   If latency-monitor-threshold is set to a value greater than 0, all events blocking the server
   for a time equal to or greater than the configured latency-monitor-threshold will be logged. Value range: 0-86400000. Default value: 0.
   **Works only on Redis.**
   + `maxclients` - (Optional, String) Set the max number of connected clients at the same time. Value range: 1000-50000. Default value: 10000.
   **Works on Redis & Memcached.**
   + `notify-keyspace-events` - (Optional, String) Redis can notify Pub or Sub clients about events happening in the key space. Default value: "Ex".
   **Works only on Redis.**
   + `repl-backlog-size` - (Optional, String) The replication backlog size in bytes for PSYNC. 
   This is the size of the buffer which accumulates slave data when slave is disconnected for some time, 
   so that when slave reconnects again, only transfer the portion of data which the slave missed. Value range: 16384-1073741824.
   Default value: 1048576. **Works only on Redis.**
   + `repl-backlog-ttl` - (Optional, String) The amount of time in seconds after the master no longer have any slaves connected
   for the master to free the replication backlog. A value of 0 means to never release the backlog. Value range: 0-604800. Default value: 3600.
   **Works only on Redis.**
   + `appendfsync` - (Optional, String) The fsync() call tells the Operating System to actually write data on disk
   instead of waiting for more data in the output buffer. Some OS will really flush data on disk,
   some other OS will just try to do it ASAP. Redis supports three different modes: 1) no: don't fsync, just let the OS flush the data when it 
   wants. 2) Faster. always: fsync after every write to the append only log. 3) Slow, Safest. everysec: fsync only one time every second. Compromise.
   Value range: "no,always,everysec". Default value: "no". **Works only on Redis.**
   + `appendonly` - (Optional, String) Configuration item for AOF persistence. Value range: "no, yes". Default value: "yes". **Works only on Redis.**
   + `slowlog-log-slower-than` - (Optional, String) The Redis Slow Log is a system to log queries that exceeded a specified execution time.
   Slowlog-log-slower-than tells Redis what is the execution time, in microseconds, to exceed in order for the command to get logged.
   Value range: 0-1000000. Default value: 10000. **Works only on Redis.**
   + `slowlog-max-len` - (Optional, String) The Redis Slow Log is a system to log queries that exceeded a specified execution time. 
   Slowlog-log-slower-than tells Redis what is the execution time, in microseconds, to exceed in order for the command to get logged.
   Value range: 0-1000. Default value: 128. **Works only on Redis.**
   + `lua-time-limit` - (Optional, String) Max execution time of a Lua script in milliseconds. Value range: 100-5000. Default value: 5000. **Works only on Redis.**
   + `repl-timeout` - (Optional, String) Replication timeout in seconds. Value range: 30-3600. Default value: 60. **Works only on Redis.**
   + `proto-max-bulk-len` - (Optional, String) Max bulk request size in bytes. Value range: 1048576-536870912. Default value: 536870912. **Works only on Redis.**
   + `master-read-only` - (Optional, String) Set redis to read only state and all write commands will fail. Value range: "yes,no". Default value: "no".
   **Works only on Redis.**
   + `client-output-buffer-slave-soft-limit` - (Optional, String) Set redis to read only state and all write commands will fail. 
   Value range: 0-1073741824. Default value: 107374182. **Works only on Redis.**
   + `client-output-buffer-slave-hard-limit` - (Optional, String) Set redis to read only state and all write commands will fail. 
   Value range: 0-1073741824. Default value: 107374182. **Works only on Redis.**
   + `client-output-buffer-limit-slave-soft-seconds` - (Optional, String) Set redis to read only state and all write commands will fail.
   Value range: 0-60. Default value: 60. **Works only on Redis.**
   + `active-expire-num` - (Optional, String) How many keys can be freed by expire cycle. Value range: 1-1000. Default value: 20.
   **Works only on Redis.**
   + `reserved-memory-percent` - (Optional, String) The percent of memory reserved for non-cache memory usage. You may want to increase 
   this parameter for nodes with read replicas, AOF enabled, etc, to reduce swap usage. Value range: 0-80. Default value: 0.
   **Works only on Memcached.**

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `configuration_parameters` - Indicates the parameter configuration defined by users based on the default parameters.

   + `name` - Indicates the parameter name.
   + `value` - Indicates the parameter value.
   + `type` - Indicates the parameter type.
   + `need_restart` - Indicates whether a restart is required.
   + `user_permission` - Indicates a user permission
