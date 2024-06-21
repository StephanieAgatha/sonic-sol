running nya pake vps aja pake command ini

```shell
wget https://go.dev/dl/go1.21.4.linux-amd64.tar.gz -O go.tar.gz
```

```shell
sudo tar -xzvf go.tar.gz -C /usr/local
```

```shell
echo export PATH=$HOME/go/bin:/usr/local/go/bin:$PATH >> ~/.profile
```

```shell
source ~/.profile
```

```shell
https://github.com/StephanieAgatha/sonic-sol.git
```

```shell
cd sonic-sol
```

paste private key in .env file
```shell
nano .env
```

```go
go run main.go
```

