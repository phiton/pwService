# password as a service - Start up an HTTP service which exposes contents of the systems passwd and groups.


pwService has two optional parameters at the commandline which allows the user to customize the location of the groups and password file of the machine. If the parameters are not present, the default location will be set to /etc/passwd and /etc/group.

# Prerequisites:
  1. GoLang installed
  2. get gorilla/mux by issuing the following command on mac: go get -u github.com/gorilla/mux

# How to run the Service:
1. Build the executable with the "go build" command while in the main directory.
2. Run the pwService executable: "./pwService". There is also the -h option: "./pwService -h" to see the help menu
3. Optional parameters can be used to configure the passwd and group file location. To configure the file locations, the command would look like "./pwService -p=<pathOfPasswdFile> -g=<pathOfGroupFile>"
      ex: ./pwService -p=/etc/myPasswdFile -g=/etc/myGroupFile
  
# How to access the service:
The service can be accessed via any webbrowser and runs on port 8081. The url will be localhost:8081.

# Available methods:
 
GET /users
Return a list of all users on the system, as defined in the /etc/passwd file.
Example Response:
[
{“name”: “root”, “uid”: 0, “gid”: 0, “comment”: “root”, “home”: “/root”,
“shell”: “/bin/bash”},
{“name”: “dwoodlins”, “uid”: 1001, “gid”: 1001, “comment”: “”, “home”:
“/home/dwoodlins”, “shell”: “/bin/false”}
]


GET
/users/query[?name=<nq>][&uid=<uq>][&gid=<gq>][&comment=<cq>][&home=<
hq>][&shell=<sq>]
Return a list of users matching all of the specified query fields. The bracket notation indicates that any of the
following query parameters may be supplied:
- name
- uid
- gid
- comment
- home
- shell
Only exact matches need to be supported.
Example Query: GET /users/query?shell=%2Fbin%2Ffalse
Example Response:
[
{“name”: “dwoodlins”, “uid”: 1001, “gid”: 1001, “comment”: “”, “home”:
“/home/dwoodlins”, “shell”: “/bin/false”}
]
  
  
GET /users/<uid>
Return a single user with <uid>. Return 404 if <uid> is not found.
Example Response:
{“name”: “dwoodlins”, “uid”: 1001, “gid”: 1001, “comment”: “”, “home”:
“/home/dwoodlins”, “shell”: “/bin/false”}
GET /users/<uid>/groups
Return all the groups for a given user.
Example Response:
[
{“name”: “docker”, “gid”: 1002, “members”: [“dwoodlins”]}
]
  
  
GET /groups
Return a list of all groups on the system, a defined by /etc/group.
Example Response:
[
{“name”: “_analyticsusers”, “gid”: 250, “members”:
[“_analyticsd’,”_networkd”,”_timed”]},
{“name”: “docker”, “gid”: 1002, “members”: []}
]


GET
/groups/query[?name=<nq>][&gid=<gq>][&member=<mq1>[&member=<mq2>][&.
..]]
Return a list of groups matching all of the specified query fields. The bracket notation indicates that any of the
following query parameters may be supplied:
- name
- gid
- member (repeated)
Any group containing all the specified members should be returned, i.e. when query members are a subset of
group members.
Example Query: GET /groups/query?member=_analyticsd&member=_networkd
Example Response:
[
{“name”: “_analyticsusers”, “gid”: 250, “members”:
[“_analyticsd’,”_networkd”,”_timed”]}
]
  
  
GET /groups/<gid>
Return a single group with <gid>. Return 404 if <gid> is not found.
Example Response:
{“name”: “docker”, “gid”: 1002, “members”: [“dwoodlins”]}
