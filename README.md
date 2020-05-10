# gitery
Go Web Service

## Deploy

### Prerequisites
1. A Linux machine with root privileges
2. SSH access to the remote machine
3. Docker installed on remote machine
4. Docker Compose installed on remote machine

### On the local machine
1. Remove the `.example` extension of `*.example` files in `/deployment` and `/configs`, setup with your own configuration.
2. Setup `deploy.sh` with remote server information.
3. run `sh deploy.sh` to build and deploy to remote server

### Use SSH to get a bash shell of remote machine
1. `cd ~/gitery`
2. `docker-compose -f service-compose.yaml build --build-arg app_env=production`
3. `docker-compose -f service-compose.yaml up -d`