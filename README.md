# auditNG: Tool for container level system call auditing made effective with selective reporting

## Introduction
auditNG suite is a dockerized open source stack with customizations of fleetmanagers, Elasticsearch with specific stored queries, TensorFlow and reporting daemon into several sinks like PagerDuty, JIRA, etc, based on learning from using several open source tools.

This tool is a service end component for the client side daemon: https://github.com/ubercoolsec/go-audit-container

System call auditing on production servers has been around for a very long time. Aggregating system call events from Linux's audit component using auditd daemon has been time tested. However, given the amount of auditd logs that get generated on a daily basis, most of which are routine, administrators go blind to typical priviledge escalation attempts like failed sudo accesses, failed multiple login attempts, unauthorized file access, etc.

When we aggregate system calls from all containers and host level nodes into a central Elasticsearch cluster, drill down into specific attributes like user id, source ip address, offending application, etc and apply machine learning, we can get a lot more insight into security events and can detect and report anomalies more effectively.

As a further step when anomalies are detected and classified into green, yellow and red and administrators are notified on high severity alerts only, it works best in getting their attention when it matters.

### Life of an audit event through auditNG
<img src="https://github.com/rhonnava/auditNG/blob/master/wiki/event_flow.png" width="600" height="300">


The stack includes:
 - Elasticsearch
 - ElastAlert with customization
 - Searching and Reporting utility
 - Tensorflow
 - Pre-processing code

Machine learning features
Features:

Time interval: Fixed
system call (open, read, write, unlink, setuid, setgid)
container_id
uid_map
exit status (Unauthorized): EACCES (-13) or EPERM (-1)
count
Score
Classification:
RED, YELLOW, GREEN


## Anomaly detection usecases
### Privilege escalation
One of the anomalies that this tool will help detect is potential privilege escalation attempts across all the monitored nodes. 


#### Privilege escalation lifecycle

<img src="https://github.com/rhonnava/auditNG/blob/master/wiki/hacking_cycle.png" width="600" height="500">


##### Reconnaissance and scanning
Privilege escalation attempts start with reconnaissance. Detection at this phase helps nipping attacks in the bud. auditNG helps drilldown to such potential privilege escalation attempts across the network on different hosts and classify the actions into red/yellow/green. Below are some of the reconnaissance attempts that are detected:

 - Listing all setuid files
 - Locating custom user accounts with some 'known default' uids. (0, 500, 501, 502, 1000, 1001, 1002, 2000, 2001, 2002)
 - Frequent calls to whoami
 - Several calls to lastlog across several hosts
 - Check to see if any hashes are stored in /etc/passwd (grep -v '^[^:]*:[x]' /etc/passwd 2>/dev/null)
 - Access to /etc/passwd, /etc/shadow, /etc/master.passwd, etc from vim, cat, nano, grep, etc
 - Looking if we can sudo on a box without supplying password (echo '' | sudo -S -l 2>/dev/null)
 - Looking for know good breakout binaries for sudo based exploits. (echo '' | sudo -S -l 2>/dev/null | grep -w 'nmap\|perl\|'awk'\|'find'\|'bash'\|'sh'\|'man'\|'more'\|'less'\|'vi'\|'emacs'\|'vim'\|'nc'\|'netcat'\|python\|ruby\|lua\|irb' | xargs -r ls -la 2>/dev/null) 
 - Searches for rhost entries
 - Looking for nfs shares and permissions (ls -la /etc/exports 2>/dev/null; cat /etc/exports 2>/dev/null)
 - Looking for credentials in /etc/fstab



##### Privilege escalation
Next, auditNG helps detect some well know privilege escalation attempts:

 - Dirty COW exploit
 - Detecting tar checkpoint command execution vilnerability by triggering an open/creat syscall with pattern "checkpoint"
 - find, vi, tar, nmap, etc commands with command execution parameters (Eg. find -exec could potentially lead to suid exploit if suid is set for find)

##### Clearing trails
Below are some typical attempts to clearing trails that will be flagged as anomalies by auditNG:

 - Write system call from editors into ~/.*_history, /root/.*_history, etc or any history -d attempts.
 - Write system calls to syslog files like /var/log/messages from any process other than syslog daemon.


##### references
 - https://www.sans.org/reading-room/whitepapers/testing/attack-defend-linux-privilege-escalation-techniques-2016-37562
 - https://github.com/dirtycow/dirtycow.github.io/blob/master/dirtyc0w.c
 - https://github.com/rebootuser/LinEnum.git
 - http://www.hackingarticles.in/4-ways-get-linux-privilege-escalation/

### East - west threats
Lateral movement of malware across hosts behind a firewall, goes undetected if there is no netflow monitoring mechanism in place. An easier way to achieve this drilling down the audit logs for connect system calls. Connect system calls are audited along with the parameters, resulting in insight on all connections made from a given host. 
Configuring allowed rules in the queries and flagging anything else as an anomaly helps in detecting such lateral spreading of malware.

### File integrity monitoring
File integrity monitoring is needed to monitor all access activities on sensitive files. It is also mandated by many regulatory requirements like PCI DSS. 

Certain sensitive files are expected to be accessed only by a whitelisted set of applications. For example, an application server, say tomcat is expected to maintain its private keys (needed for TLS handshakes) in a .pem file on the file system. However, only the application server is expected to access it and no other application. Any anomaly from the whitelisted applications, with respect to such file access can be flagged as an anomaly.

### ### Container breakout detection
Detection of a file access by a process with a container tag, having broken out of the namespace.

