# otc

**otc** -- out to cloud, has for goal to simplify storing files, or making database backups periodically.

## Installation

For now, you need the go toolchain to install otc

 ```sh
go install github.com/Karitham/otc@latest
 ```

## Examples

ping a google every 5 second and send the result in a discord file

```sh
otc cron --schedule 5s cmd --command ping --args '-c 1 google.com' discord --file 'ping-results.txt' --url $WEBHOOK_URL
```
