box.cfg{
    listen = 3301,
    pid_file = nil,
    background = true,
    log_level = 5,
    memtx_memory = 512 * 1024 * 1024,
    memtx_max_tuple_size = 8 * 1024 * 1024
}

box.schema.user.create('user', {password='password', if_not_exists=true})
box.schema.user.grant('user', 'read,write', 'universe')

s = box.schema.space.create('users')
s:format({
         {name = 'id', type = 'unsigned'},
         {name = 'first_name', type = 'string'},
         {name = 'last_name', type = 'string'},
         {name = 'age', type = 'unsigned'},
         {name = 'sex', type = 'unsigned'},
         {name = 'interests', type = 'string'},
         {name = 'city', type = 'string'}
         })

s:create_index('primary', {
        type = 'tree',
        parts = {'id'}
        })

s:create_index(
        'secondary',
        {
            type = 'tree',
            unique = false,
            parts = {'first_name', 'last_name'}
        }
)


s:insert{10000000000, 'first_name', 'last_name', 25, 1, 'fdsafdasfds','Moscow'}
