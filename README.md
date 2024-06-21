use vps for better experience

use this command if u don't have go on your machine.if u already installed go, just paste ur private key to .env file then run code

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

clone repo
```shell
git clone https://github.com/StephanieAgatha/sonic-sol.git
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

If you want to monitoring your total transaction on sonic, just paste your Authorization value from sonic server.you can get it using devtools and go to network tab > choose fetch/XHR then search "daily". you will see your jwt token there.

![image](https://github.com/StephanieAgatha/sonic-sol/assets/62786809/88c14ef1-838d-4ba6-b2c2-7238487c7c77)

copy and input your jwt token to bot.

and if you see "failed to send transaction xxxx, timeout with tx signature about (looks like this picture) 
![image](https://github.com/StephanieAgatha/sonic-sol/assets/62786809/83fc47dd-6bb5-477c-a8d4-b196c9e3f255)

it's fine,let bot run
