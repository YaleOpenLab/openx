node {
  docker.image('golang').inside {

		stage ('Setup env vars') {
			sh 'useradd -ms /bin/bash jenkins'
			sh 'sudo su jenkins'
			sh 'echo "GOROOT=/usr/local/go" >> ~/.profile'
			sh 'echo "GOPATH=$HOME" >> ~/.profile'
			sh 'echo "PATH=$GOPATH/bin:$GOROOT/bin:$PATH" >> ~/.profile'
			sh 'export GOCACHE="/tmp/.cache"'
			sh 'export XDG_CACHE_HOME="/tmp/.cache"'
		}
		stage('Print go version') {
    	sh 'go version'
  	}

		stage('Get openx package') {
			sh 'go get -v github.com/YaleOpenLab/openx'
		}
  }
}
