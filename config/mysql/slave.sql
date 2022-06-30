CHANGE MASTER TO MASTER_HOST='db-master', MASTER_USER='repl', MASTER_PASSWORD='slavepass';

# INSTALL PLUGIN rpl_semi_sync_slave SONAME 'semisync_slave.so';
# SET GLOBAL rpl_semi_sync_slave_enabled=0;
