pipeline {

	environment{
		WORKSPACE = '/home/jenkins'
	}

	stages {
		stage ('Test') {
			steps {
 				sh "printenv | sort"
				sh 'sudo apt-get -y upgrade'
				sh 'wget https://dl.google.com/go/go1.12.4.linux-amd64.tar.gz'
				sh 'sudo tar -xvf go1.12.4.linux-amd64.tar.gz'
				sh 'sudo mv go /usr/local'
				sh 'export GOROOT="/usr/local/go"'
				sh 'export GOPATH="$HOME"'
				sh 'export PATH="$GOPATH/bin:$GOROOT/bin:$PATH"'
				sh 'go version'
				sh 'go get -v github.com/YaleOpenLab/openx'
			}
		}
	}
}
