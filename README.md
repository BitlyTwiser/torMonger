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
- For this application to function, open does need to have Tor installed. 
- for linux (debian/ubuntu, or any distro using apt) ```sudo apt install tor```
- When starting the application, you can select the given port to match what exists within your torrc configuration file.
- ```sudo vim /etc/tor/torrc```
- Then run tor from the cli (or tor browser) and connect
- ```tor```
- You can now run torMonger using the port designated in the torrc file.
- Validate ports and the service using netstat.
- ```netstat -tulpn```
- Note: The port might be set to 9050, thus you will have to account for this when running the application.

## Usage:
- The torMonger application can be installed by compiling the Go binary and running the install script.
- With the use of Go modules you can perform the installation with ```go install``` then compile the binary.
- ```sudo ./install.sh```
- This application was tested and ran on Linux. (Tails OS was used for the runtime environment)
- After installation, one can run the command line tool.
- ```torMongert -url http://zqktlwiuavvvqqt4ybvgvi7tyo4hjl5xgfuvpdf6otjiycgwqbym2qad.onion/wiki/ -t 1```
- The above onion link is associated with "The Hidden Wiki".
- The ```-threads``` flag spawns the selected number of processes to ingest data.

##Docker Images:
- All data is initially stored within the mongoDB instance.
- Errors logs are within the logs table
- Simply just run docker-compose against the provided yaml file to spawn all necessary containers in the environment.
- ```docker-compose up -d --build```
- The following command with daemonize the docker processes and rebuild the containers images from scratch


