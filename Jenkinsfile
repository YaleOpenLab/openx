pipeline {

	agent none

	environment{
		WORKSPACE = '/home/jenkins'
	}

	stages {
		stage ('Test') {
			agent { label 'master'}
			steps {
				sh 'sudo rm -rf /usr/local/go'
 				sh "printenv | sort"
				sh 'wget https://dl.google.com/go/go1.12.4.linux-amd64.tar.gz'
				sh 'tar -xvf go1.12.4.linux-amd64.tar.gz'
				sh 'sudo mv go /usr/local'
				sh 'echo "GOROOT=/usr/local/go" >> ~/.profile'
				sh 'echo "GOPATH=$HOME" >> ~/.profile'
				sh 'echo "PATH=$GOPATH/bin:$GOROOT/bin:$PATH" >> ~/.profile'
				sh '/usr/local/go/bin/go version'
				sh '/usr/local/go/bin/go get -v github.com/YaleOpenLab/openx'
			}
		}
	}
}
