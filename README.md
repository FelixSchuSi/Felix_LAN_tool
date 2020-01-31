# Felix_LAN_tool
This tool enables you to transfer data within a local network over a web interface.
There are plenty other, probably better methods of transfering data within a LAN. The intention behind this project was to get first hand development experience with NodeJS and Handlebars. On the plus side, this method is easy to use for end users.
You can access the interface in your browser by using the local IP address of the host as the URL.
This tool runs on Port 3000 by default. For best performance I reccomend to set up a reverse proxy e.g. nginx, apache.

## run
`git clone https://github.com/FelixSchuSi/Felix_LAN_tool.git`

`cd Felix_LAN_tool`

`npm i`

Place your data in the 'files-to-be-transfered' directory.

`npm run`

access web interface over local ip address:port

## determine local IP address
The local ip address of your machine can be identified with `ipconfig` on Windows and `ifconfig` on Linux/Mac.

