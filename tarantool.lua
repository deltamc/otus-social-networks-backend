box.cfg{
    listen = 3301,
     background = true,
       log = '1.log',
       pid_file = '1.pid'
}

box.schema.user.create('user', {password='password', if_not_exists = true})
box.schema.user.grant('user', 'read,write,execute', 'universe', nil, {if_not_exists=true})

box.space.users:drop()

s = box.schema.space.create('users', {id=9, field_count=7, if_not_exists=true})
s:format({
         {name = 'id', type = 'integer'},
         {name = 'first_name', type = 'string'},
         {name = 'last_name', type = 'string'},
         {name = 'age', type = 'integer'},
         {name = 'sex', type = 'integer'},
         {name = 'interests', type = 'string'},
         {name = 'city', type = 'string'}
         })
s:create_index('primary', {
        type = 'tree',
        parts = {'id'},
        if_not_exists=true
        })

s:create_index(
        'secondary',
        {
            type = 'tree',
            unique = false,
            parts = {'first_name', 'last_name'},
            if_not_exists=true
        }

)
s:create_index(
        'fn',
        {
            type = 'tree',
            unique = false,
            parts = {'first_name'},
            if_not_exists=true
        }
)
s:create_index(
        'ln',
        {
            type = 'tree',
            unique = false,
            parts = {'last_name'},
            if_not_exists=true
        }
)
--
--
s:insert{10000000000, 'first_name', 'last_name', 25, 1, 'fdsafdasfds','Moscow'}
local sp = box.space._space.index.name:select{ 'users' }
print(sp[1])

function get_users(first_name, last_name)
    if first_name ~= "" and last_name ~= "" then
        return box.space.users.index.secondary:select({first_name, last_name})
    else
        if first_name ~= "" and last_name == "" then
            return box.space.users.index.fn:select({first_name})
        else
            if first_name == "" and last_name ~= "" then
                return box.space.users.index.ln:select({last_name})
            else
                return box.space.users:select({})
            end
       end
    end
end