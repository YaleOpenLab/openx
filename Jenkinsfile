pipeline {

	agent {label 'master'}

	environment{
		WORKSPACE = '/home/jenkins'
	}

	stages {
		stage ('Install Golang') {
			steps {
				sh 'go get -v github.com/YaleOpenLab/openx'
				sh 'cd ~/go/src/github.com/YaleOpenLab/openx'
				sh '/usr/local/go/bin/go get ./...'
				sh '/usr/local/go/bin/go build -v ./...'
			}
		}
	}
}
