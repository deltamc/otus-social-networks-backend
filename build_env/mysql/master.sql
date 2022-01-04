CREATE USER repl@'%' IDENTIFIED WITH mysql_native_password BY 'slavepass';
GRANT REPLICATION SLAVE ON *.* TO repl@'%';

# INSTALL PLUGIN rpl_semi_sync_master SONAME 'semisync_master.so';

# set global rpl_semi_sync_master_enabled=0;
