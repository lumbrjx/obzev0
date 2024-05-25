<h1 align="center">TON-618</h1>

<p align="center">
 PoC of a chaos testing tool for tcp connections written in go
</p>

<p align="center">
  <img src="./assets/tn.jpg" />
</p>

## How to configure this tool
in the same dir of ton618 binary create a ``tonConf.yaml`` file with the following attr's:

```yaml
qDelay: 0  // here you can set the delay for requests and responses
  resDelay: 0
server:
  port: "7090" // tcp server port
```
then you can start it via ``./ton618``

## How to use it

start your http server then use curl or nc to send request via the proxy server 
*you can use the http server examg go run httpServer.go*

