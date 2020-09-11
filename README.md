# cloud_spanner_api



## 1. Install golang and prepare the environment

```
sudo apt-get update
sudo apt-get -y upgrade

sudo apt-get install -y git
sudo apt-get install -y wget

wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz
tar -xvf go1.13.3.linux-amd64.tar.gz
sudo mv go /usr/local

export GOROOT=/usr/local/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

go get cloud.google.com/go/spanner
```

## 2. Create the API credentials

Visit the credentials page: https://console.cloud.google.com/apis/credentials.

Click on the ```Create service account``` button and follow the instructions.
Save the resulting json file somewhere on the filesystem, e.g. ```/home/user/token.json```

## 3. Enable Cloud Resource Manager API for your projects

In order to access a project with API, one needs to enable Cloud Resource Manager API for that project.

Visit https://console.cloud.google.com/apis/library/cloudresourcemanager.googleapis.com, press the ```Enable``` button and choose your project. 
