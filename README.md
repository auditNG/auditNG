# auditNG: Tool for container level system call auditing made effective with selective reporting

This tool is a service end component for the client side daemon: https://github.com/ubercoolsec/go-audit-container

System call auditing on production servers has been around for a very long time. Aggregating system call events from Linux's audit component using auditd daemon has been time tested. However, given the amount of auditd logs that get generated on a daily basis, most of which are routine, administrators go blind to typical priviledge escalation attempts like failed sudo accesses, failed multiple login attempts, unauthorized file access, etc.

When we aggregate system calls from all containers and host level nodes into a central Elasticsearch cluster, drill down into specific attributes like user id, source ip address, offending application, etc and apply machine learning, we can get a lot more insight into security events and can detect and report anomalies more effectively.

As a further step when anomalies are detected and classified into green, yellow and red and administrators are notified on high severity alerts only, it works best in getting their attention when it matters.

auditNG suite is a dockerized open source stack with customizations of fleetmanagers, Elasticsearch with specific stored queries, TensorFlow and reporting daemon into several sinks like PagerDuty, JIRA, etc, based on learning from using several open source tools.

The stack includes:
 - Elasticsearch
 - ElastAlert with customization
 - Searching and Reporting utility
 - Tensorflow
 - Pre-processing code
