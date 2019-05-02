# B2Backups

A golang client for Backblaze's B2 storage. This is a work in progress, Currently only uploads file changes.

See Issues for Work in Progress

## Usage

Some simple examples to get you started. Errors are ommitted for brevity. 


Initialize a folder
~~~
b2backup -init=<path>
~~~

Startup
~~~
./b2backup
~~~

AutoStart
~~~
systemctl link /path/to/b2backup.service
systemctl enable b2backup
~~~