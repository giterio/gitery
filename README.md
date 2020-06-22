# gitery
Gitery Web Service

## Deploy

### Prerequisites
1. An account of Linux machine with sudo privileges
2. Authorited SSH access to the remote machine
3. Docker installed on remote machine
4. Docker Compose installed on remote machine

### On the local machine
1. Remove the `.example` extension of `*.example` files in `/deployment` and `/configs`, setup with your own configuration.
2. Setup `build.sh` with your remote server SSH information.
3. run `sh build.sh` to make a production build and upload `/deployment` and `/bin` to remote server via ssh connection.

### Use SSH to get a bash shell of remote machine
1. `cd ~/gitery`
2. run `sh deploy.sh` to start the service on a Nginx server