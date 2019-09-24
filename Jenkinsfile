node {
    withEnv(["WORKSPACE=/home/jenkins", "GOPATH=${env.WORKSPACE}/go", "GOROOT=${root}", "GOBIN=${root}/bin", "PATH+GO=${root}/bin"]) {
			docker.image('golang').inside {
				sh 'echo "$GOPATH"'
	    	sh 'go get -v github.com/YaleOpenLab/openx'
		  }
    }
}
