###Description
This application checkes docker containers 
and when any containers stopped, it will 
send message to mattermost
###Project structure
1. *__config__* - this package parses config.toml
2. *__database__* - contains CRUD for leveldb. Consist of 2 tables: 
    * first for saving key-pair container_name - chat_name
    * second for saving key-pair user_name - user_pass
3. *__example_config__* - contains example of config file
4. *__front__* - web interface for application
5. *__mattermost__* - workes with mattermost, 
used its SDK. First of all it creates object of 
mattermost client and get JWT token from mattermost. 
Secondly it waits when chanel for transfering message 
will send something. Thirdly sends message from go-chanel 
to mattermost-chanel.
6. *__portainer__* - checkes of docker containers pool, 
if one on them stoped, sends its name and logs by the 
chanel to mattermost
7. *__rest__* - Create JWT token for auth. 
Contains handlers, which provides working with CRUD.

###Run and install

Project contains makefile. So ```make run``` will running 
this app, but at start it will copy __config.toml__ 
from *example_config* to *$HOME/.portainer_plugin*. 
```make install``` will create binary file and also 
will copied config file.

###Config

In config file contains settings for each package. 
Address and port for connection to portainer and mattermost, 
Email/login and password of accounts through wich application 
will work with portainer/mattermost, and contains logger 
settings, path to storage, port for rest, checking 
containers interval and and amount strings of  logs 
from stoped container