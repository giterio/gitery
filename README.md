# gitery
Gitery Web Service

## Deploy direction

### Prerequisites
1. An account of Linux machine with sudo privileges
2. Authorited SSH access to the remote machine
3. Docker installed on remote machine
4. Docker Compose installed on remote machine

### On the local machine
1. Remove the `.sample` extension of `*.sample` files in `/deployment`, `/configs` and root directory, setup with your own configuration.
2. Run `sh build.sh` to make a production build and upload `/deployment` and `/bin` to remote server via ssh connection.

### Use SSH to get a bash shell of remote machine
1. `cd ~/gitery`
2. run `sh deploy.sh` to start the service on a Nginx server

### Setup the database on remote server
1. Suppose the remote server IP address is `xxx.xxx.xxx.xxx`, please visit `xxx.xxx.xxx.xxx:5050` in your browser.
2. Login the pgAdmin 4 (PostgreSQL GUI) with the account infomation configured in `/deployment/service-compose.yaml`.
3. Setup postgreSQL server with the configuration infomation in `/deployment/service-compose.yaml`. Please refer to [pgAdmin 4 document](https://www.pgadmin.org/) if you have any problem.
4. Initialize the database with the SQL clause in `/internal/database/database.sql`.

#### After finishing all the steps mentioned above, the gitery service should be ready then.