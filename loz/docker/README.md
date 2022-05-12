# Readme
### Docker image too big for git 
### Saved on box : https://ibm.ent.box.com/folder/144242217897
#### (1) TFS image: image_csl_ccf_os_loz_tfs_v1.tar (PORT 8500 for REST API ; PORT 8501 for gRPC)
##### command to start container: docker run -d -t  -e MODEL_NAME=model --name ccf_os_loz_tfs_v1 csl/ccf_os_loz_tfs:v1 

#### (2) REST API image: image_csl_ccf_os_loz_api_svr_v1.tar (PORT 8080 for REST API)
##### command to start container: docker run -d -t -p 8080:8080  -e REDISADD="tent-redis:6379"  -e TFSADD="ccf_os_loz_tfs_v1:8501"  --name ccf_os_loz_api_svr_v1 --link tent-redis --link  ccf_os_loz_tfs_v1 csl/ccf_os_loz_api_svr:v1
#####  As Redis and TFS container on the same image, using link to contact REST API server with Redis and TFS
#####  "tent-redis" is your Redis server container name 
#####  "ccf_os_loz_tfs_v1" is your TFS container name
#####  REDISADD is the your Redis server container name  plus port number
#####  TFSADD is the your TFS server container name  plus port number

