pipeline {
    agent none
    environment {
        NODENAME = "cplaton"
    }
    stages {
        stage("Build") {
            environment {
                GOCACHE = "${WORKSPACE}"
            }

            agent {
                docker {
                    image "platon-builder:1.0.1"
                }
            }

            steps {
                sh "chmod u+x ./build/*.sh"
                sh "go mod tidy"
                sh "./build/build_deps.sh"
                sh "make clean && make platon"
                sh "mv ./build/bin/platon ./build/bin/${NODENAME}"
                archiveArtifacts artifacts: "build/bin/${NODENAME}"
            }
        }
    }
}