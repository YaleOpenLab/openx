node {
    try{
        ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
            withEnv(["GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"]) {
                env.PATH="${GOPATH}/bin:$PATH"

                stage('Checkout'){
                    echo 'Checking out SCM'
                    checkout scm
                }

                stage('Build'){
                    echo 'Building Executable'
										sh """cd $GOPATH && go build"""
                }
            }
        }
    }catch (e) {
        // If there was an exception thrown, the build failed
        currentBuild.result = "FAILED"
    }
}
