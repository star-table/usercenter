#!/usr/bin/env groovy
//def BRANCH_NAME = git_branch.tokenize('/')[-2]+"/"+git_branch.tokenize('/')[-1]
def BRANCH_NAME = git_branch.replace("origin/","")
pipeline{
	agent {
		label 'jenkins-slave'
	}
    options {
        disableConcurrentBuilds()
        skipDefaultCheckout()
        timeout(time: 1, unit: 'HOURS')
        timestamps()
    }
	environment{
		GIT_URL = ""
		CREDENTIALS_ID = ""
		app = "lesscode-usercenter"
	}
	stages {
        stage('print vars') {
			steps {
				echo "${BRANCH_NAME}"
				echo "${env.env}"
				}
        }
		stage('checkout') {
			steps {
				git branch: "${BRANCH_NAME}", credentialsId: "${CREDENTIALS_ID}", url: "${GIT_URL}"
			}
		}
		stage('package') {
			steps {
				sh '''
				CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o usercenter main.go
				'''
			}
		}
		stage('docker build') {
			steps {
				sh '''
                echo "build docker image for lesscode-usercenter"
				if [ ${env}"bjx" = "devbjx" -o ${env}"bjx" = "testbjx" -o ${env}"bjx" = "k8stestbjx" -o ${env}"bjx" = "crmbjx" -o ${env}"bjx" = "crm2bjx" -o ${env}"bjx" = "k8s-testbjx" -o ${env}"bjx" = "fusebjx" -o ${env}"bjx" = "graybjx" ];then
				if [ ${env}"bjx" = "k8s-testbjx" ];then
				sed -i 's/15001/15003/g' Dockerfile
				fi
				docker build -t registry-vpc.cn-shanghai.aliyuncs.com/polaris-team/lesscode-usercenter:v1.${BUILD_ID} .
				else
                docker build -f Dockerfile -t harbor.bd88888.online/bigdata/lesscode-usercenter:v1.${BUILD_ID} .
				fi
                '''
			}
		}
		stage('docker push') {
			steps {
				sh '''
				if [ ${env}"bjx" = "devbjx" -o ${env}"bjx" = "testbjx" -o ${env}"bjx" = "k8stestbjx" -o ${env}"bjx" = "crmbjx" -o ${env}"bjx" = "crm2bjx" -o ${env}"bjx" = "k8s-testbjx" -o ${env}"bjx" = "fusebjx" -o ${env}"bjx" = "graybjx" ];then
                docker push registry-vpc.cn-shanghai.aliyuncs.com/polaris-team/lesscode-usercenter:v1.${BUILD_ID}
				else
				docker push harbor.bd88888.online/bigdata/lesscode-usercenter:v1.${BUILD_ID}
				fi
                '''
			}
		}
		stage('deploy') {
			steps {
				sh '''
				if [ ${env}bjx = "devbjx" ];then
				ssh media@172.19.84.54 -p 19222 "cd /data/app/lesscode-usercenter && ./upgrade.sh v1.${BUILD_ID}"
				elif [ ${env}bjx = "k8s-testbjx" ];then
				ssh media@172.19.84.54 -p 11222 "kubectl -ntest set image deployment/lesscode-usercenter lesscode-usercenter=registry.cn-shanghai.aliyuncs.com/polaris-team/lesscode-usercenter:v1.${BUILD_ID} --record"
				elif [ ${env}bjx = "testbjx" ];then
				ssh media@172.19.132.103 "cd /data/app/lesscode-usercenter && ./upgrade.sh v1.${BUILD_ID}"
				elif [ ${env}bjx = "crmbjx" ];then
				ssh media@172.19.84.54 -p 17222 "cd /data/app/lesscode-usercenter && ./upgrade.sh v1.${BUILD_ID}"
				elif [ ${env}bjx = "crm2bjx" ];then
				ssh media@172.19.71.36 "cd /data/app/lesscode/lesscode-usercenter && ./upgrade.sh v1.${BUILD_ID}"
				elif [ ${env}bjx = "test2bjx" ];then
				ssh media@172.19.132.102 "cd /data/app/lesscode/lesscode-usercenter && ./upgrade.sh v1.${BUILD_ID}"
				elif [ ${env}bjx = "k8stestbjx" ];then
				kubectl -ntest set image deployment/lesscode-usercenter lesscode-usercenter=registry.cn-shanghai.aliyuncs.com/polaris-team/lesscode-usercenter:v1.${BUILD_ID} --record
				elif [ ${env}bjx = "fusebjx" ];then
				kubectl -nfuse set image deployment/lesscode-usercenter lesscode-usercenter=registry.cn-shanghai.aliyuncs.com/polaris-team/lesscode-usercenter:v1.${BUILD_ID} --record
				elif [ ${env}bjx = "graybjx" ];then
				kubectl -ngray set image deployment/lesscode-usercenter lesscode-usercenter=registry.cn-shanghai.aliyuncs.com/polaris-team/lesscode-usercenter:v1.${BUILD_ID} --record
				elif [ ${env}bjx = "prodbjx" ];then
				echo "./upgrade.sh v1.${BUILD_ID}"
				elif [ ${env}bjx = "stagbjx" ];then
				echo "./upgrade.sh v1.${BUILD_ID}"
				ssh media@18.167.142.236 "./lesscode-deploy.sh 46.137.228.152 lesscode-usercenter v1.${BUILD_ID}"
				elif [ ${env}bjx = "stag_twbjx" ];then
				echo "./upgrade.sh v1.${BUILD_ID}"
				ssh media@18.167.142.236 "./polaris-deploy.sh 46.137.228.152 polaris-usercenter v1.${BUILD_ID}"
				else
				echo "skip deploy"
				fi
                '''
			}
		}
	}
}
