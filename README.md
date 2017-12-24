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

The output of go-audit-container is as follows:
Success:
Dec 20 09:03:14 idc-dost-bm03.dev.walmart.com go-audit[11358]: {"sequence":861471,"timestamp":"1513740794.529","messages":[{"type":1300,"data":"arch=c000003e syscall=59 success=yes exit=0 a0=ec87d0 a1=ec64a0 a2=ed5490 a3=7ffccb076780 items=2 ppid=11461 pid=11616 auid=1006 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=pts2 ses=8137 comm=\"grep\" exe=\"/usr/bin/grep\" subj=unconfined_u:unconfined_r:unconfined_t:s0-s0:c0.c1023 key=\"user_commands\""},{"type":1309,"data":"argc=3 a0=\"grep\" a1=\"--color=auto\" a2=\"go-audit\""},{"type":1307,"data":" cwd=\"/root\""},{"type":1302,"data":"item=0 name=\"/bin/grep\" inode=100667447 dev=fd:00 mode=0100755 ouid=0 ogid=0 rdev=00:00 obj=system_u:object_r:bin_t:s0 objtype=NORMAL"},{"type":1302,"data":"item=1 name=\"/lib64/ld-linux-x86-64.so.2\" inode=43653168 dev=fd:00 mode=0100755 ouid=0 ogid=0 rdev=00:00 obj=system_u:object_r:ld_so_t:s0 objtype=NORMAL"}],"uid_map":{"0":"root","1006":"rhonnav"},"container_id":0}

Dec 20 09:57:10 idc-dost-bm03.dev.walmart.com go-audit[11358]: {"sequence":862017,"timestamp":"1513744030.477","messages":[{"type":1300,"data":"arch=c000003e syscall=59 success=yes exit=0 a0=168fb00 a1=17988f0 a2=16814d0 a3=7ffed1fbfa50 items=2 ppid=11401 pid=12276 auid=1006 uid=1006 gid=1006 euid=1006 suid=1006 fsuid=1006 egid=1006 sgid=1006 fsgid=1006 tty=pts2 ses=8137 comm=\"touch\" exe=\"/usr/bin/touch\" subj=unconfined_u:unconfined_r:unconfined_t:s0-s0:c0.c1023 key=\"user_commands\""},{"type":1309,"data":"argc=2 a0=\"touch\" a1=\"/root/abc\""},{"type":1307,"data":" cwd=\"/home/rhonnav\""},{"type":1302,"data":"item=0 name=\"/usr/bin/touch\" inode=100692532 dev=fd:00 mode=0100755 ouid=0 ogid=0 rdev=00:00 obj=system_u:object_r:bin_t:s0 objtype=NORMAL"},{"type":1302,"data":"item=1 name=\"/lib64/ld-linux-x86-64.so.2\" inode=43653168 dev=fd:00 mode=0100755 ouid=0 ogid=0 rdev=00:00 obj=system_u:object_r:ld_so_t:s0 objtype=NORMAL"}],"uid_map":{"0":"root","1006":"rhonnav"},"container_id":0}

Failure:

Dec 20 09:57:10 idc-dost-bm03.dev.walmart.com go-audit[11358]: {"sequence":862018,"timestamp":"1513744030.480","messages":[{"type":1300,"data":"arch=c000003e syscall=2 success=no exit=-13 a0=7fff6454c6c2 a1=941 a2=1b6 a3=7fff6454b090 items=1 ppid=11401 pid=12276 auid=1006 uid=1006 gid=1006 euid=1006 suid=1006 fsuid=1006 egid=1006 sgid=1006 fsgid=1006 tty=pts2 ses=8137 comm=\"touch\" exe=\"/usr/bin/touch\" subj=unconfined_u:unconfined_r:unconfined_t:s0-s0:c0.c1023 key=\"access\""},{"type":1307,"data":" cwd=\"/home/rhonnav\""},{"type":1302,"data":"item=0 name=\"/root/abc\" objtype=UNKNOWN"}],"uid_map":{"1006":"rhonnav"},"container_id":0}

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

Identifying privilege escalation through the different phases

![alt text](https://github.com/rhonnava/auditNG/blob/master/wiki/hacking_cycle.png?raw=true =100X100)
