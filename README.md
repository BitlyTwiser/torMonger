```
    ███      ▄██████▄     ▄████████   ▄▄▄▄███▄▄▄▄    ▄██████▄  ███▄▄▄▄      ▄██████▄     ▄████████    ▄████████ 
▀█████████▄ ███    ███   ███    ███ ▄██▀▀▀███▀▀▀██▄ ███    ███ ███▀▀▀██▄   ███    ███   ███    ███   ███    ███ 
   ▀███▀▀██ ███    ███   ███    ███ ███   ███   ███ ███    ███ ███   ███   ███    █▀    ███    █▀    ███    ███ 
    ███   ▀ ███    ███  ▄███▄▄▄▄██▀ ███   ███   ███ ███    ███ ███   ███  ▄███         ▄███▄▄▄      ▄███▄▄▄▄██▀ 
    ███     ███    ███ ▀▀███▀▀▀▀▀   ███   ███   ███ ███    ███ ███   ███ ▀▀███ ████▄  ▀▀███▀▀▀     ▀▀███▀▀▀▀▀   
    ███     ███    ███ ▀███████████ ███   ███   ███ ███    ███ ███   ███   ███    ███   ███    █▄  ▀███████████ 
    ███     ███    ███   ███    ███ ███   ███   ███ ███    ███ ███   ███   ███    ███   ███    ███   ███    ███ 
   ▄████▀    ▀██████▀    ███    ███  ▀█   ███   █▀   ▀██████▀   ▀█   █▀    ████████▀    ██████████   ███    ███ 
                         ███    ███                                                                  ███    ███ 
// Mainteiner: BitlyTwiser
// License: MIT
```

# torMonger
- A recursive Tor network crawler

## Tor:
- For this application to function, open does need to have Tor installed and utilize the PSQL instance built via the docker-compose file. 
- for linux (debian/ubuntu, or any distro using apt) ```sudo apt install tor```
- When starting the application, you can select the given port to match what exists within your torrc configuration file.
- ```sudo vim /etc/tor/torrc```
- Then run tor from the cli (or tor browser) and connect
- ```tor```
- You can now run torMonger using the port designated in the torrc file.
- Validate ports and the service using netstat.
- ```netstat -tulpn```
- Note: If you execute ```tor``` via the cli, the application initiates with port 9050. This is the default port used with tormonger, however, if you execute tor via the tor web browser, 9150 is used.

## Usage:
- The torMonger application can be installed by compiling the Go binary and running the installation script.
- With the use of Go modules you can perform the installation with ```go install``` then compile the binary.
- ```sudo ./install.sh```
- This application was tested and ran on Linux. 
- After installation, one can run the command line tool.
- ```torMonger -url http://zqktlwiuavvvqqt4ybvgvi7tyo4hjl5xgfuvpdf6otjiycgwqbym2qad.onion/wiki/ -t 1```
- The above onion link is associated with "The Hidden Wiki".
- The ```-threads``` flag spawns the selected number of processes to ingest data.

- When the application is executed one would expect to see the following output:
[Output](./images/running_example.png)

#### Env Vars:
- The application takes advantage of a .env file for loading database credentials and values.
- The .env is included within the project and can be altered to your liking. 
- Just match the .env credentials with that within the docker-compose file to connect to the database.

## Docker Images:
- All data is initially stored within the PSQL database instance, the install.sh script will create a directory ```/var/tmp/tormonger_data```.
- Errors logs are within the logs table.
- Simply just run docker-compose against the provided yaml file to spawn all necessary containers in the environment.
- ```docker-compose up -d --build```
- The following command with daemonize the docker processes and rebuild the containers images from scratch
- The docker image pulls in the ```create_tables.sql``` file which is executed when the image is built. this will create all the tables and foreign key associations for you.

## Data stroage and Example data:
- The primary column is the tormonger_data column, which stores the base link url and the base64 encoded link hash.
[tormonger data column example](./images/tormonger_data_column.png)
- When the application is running succesfully, you will obtain snapshots of the html data from each webpage within the html_data table.
[html data table](./images/html_table_example.png)
- The html is stored in the html_data column, which is direct html that can be pulled out and viewed.
[html data example](./images/html_example.png)
- subdirectories can be seen in the ```tormonger_data_sub_directories``` table.
[subdirectories](./images/sub_dirs.png)
- you can search for subdirs based off of the id for the link you desire and search html data either by subdirectory or by main link id.
- index's were used on the primary elements of each column for faster lookup times.



