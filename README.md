# PSGEN

`psgen` is CLI tool for local password management. You can generate passwords and optionally store them in a local SQLite database. The password db can be exported to a CSV file. Additionally a CSV file with passwords can be imported to the db. Stored passwords are encrypted using AES encryption.

### Config
psgen configuration and the sqlite db are both stored within ~/.psgen folder.\
Configuration example:
```
{
    "enc_key":"82sdnT0W1axo7GdhhmeFYwLsMkLXhoqKJRamkenWHRU=",
    "execution_timeout":5,
    "db_path":"/home_path/.psgen/psgen.db",
    "logs_path":"/home_path/.psgen/logs"
}
```
enc_key: a base64 string representing a sequence of 32 bytes\
execution_timeout: max number of second, that a db operation can take\
db_path: path of the sqlite db\
logs_path: errors are stored in log files, that resides in /home_path/.psgen/logs folder
### Example
```bash
$ go run main.go gen -s  -d -ln=25
password:  w0+,@,+i>29-f34^y^810wuq`
Store password? Y[yes] yes
Give password key: password-key
Password successfully generated%
```

### Usage
```bash
$ go run main.go --help 
Usage: psgen <command> -[-]<flags>
Commands:
gen             generates a password
get             retrieves a password from the local sqlite db and prints it out
export          exports the stored passwords from the local sqlite db to an csv file
import          imports passwords from a csv file into the local sqlite db
help            show help

Use psgen <command> -h or --help for more information about a command.
```

### Run
Execute using go command
```bash
$ go run main.go <command> -[-]<flags>
```

Execute binary:
1) Build the binaries for darwin and linux OS
```bash
$ make
Building binary....
GOARCH=arm64 GOOS=darwin go build -o ./builds/darwin/psgen main.go
GOARCH=amd64 GOOS=linux go build -o ./builds/linux/psgen main.go
```
2) run the binary
```bash
$ cd builds
$ ./psgen <command> -[-]<flags>
```
