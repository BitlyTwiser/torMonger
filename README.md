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
- The environment is initialized with several docker containers as well.
- Elasticsearch is used for reduced lookup times when querying data from the frontend.
- All data is initially stored within the mongoDB instance and synced with Elasticsearch on 5 minute intervals via a job runner.
- Simply just run docker-compose against the provided yaml file to spawn all necessary containers in the environment.
- ```docker-compose up -d --build```
- The following command with daemonize the docker processes and rebuild the containers images from scratch


