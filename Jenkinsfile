node {
  docker.image('golang').inside {
    withEnv(["WORKSPACE=/home/jenkins", "GOPATH=${env.WORKSPACE}/go", "GOROOT=${root}", "GOBIN=${root}/bin", "PATH+GO=${root}/bin"]) {
    	sh 'echo "$GOPATH"'
    	sh 'go get -v github.com/YaleOpenLab/openx'
    }
  }
}
