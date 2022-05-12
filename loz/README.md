# Readme
CCF OS LoZ backend REST API Server files
## 1)  api_svr_rc:
    This is an executable file. 
    It will start a REST API listen on por 8080
    Only one REST API 
    Input is a transaction index (which will used to pick transaction from test_220_100k_os.csv)
    Output is the predication value):
    URL http://IP:8080/ccf_inference
    Method POST
### HTTP Request sample:    
    Request JSON body is {"tx_index": 1105755 } 
    The 1105755 here is index of transaction in file test_220_100k.indices    
### HTTP Response sample:
{
    "Elapsed    TFS": "5.217348ms",
    "Elapsed  Redis": "333.596µs",
    "Elapsed Mapper": "158.089µs",
    "Predict  Value": "7.18149465E-07",
    "TX       Index": "1105755"
}

### Dependence:
#### i)	Redis: localhost:6367
#### ii) TensorFlow Severing : http://localhost:8501/v1/models/model:predict
#### iii) mapper.so
    
## 2) mapper.so
  This is the library used for mapping.
### Input sample:
{"Merchant_State":"missing_value","Zip":0,"Merchant_Name":"-6160036380778658394","Merchant_City":" ONLINE","MCC":4121,"Use_Chip":"Online Transaction","Errors":"missing_value","Year":2010,"Month":2,"Day":19,"Time":"21:50","Amount":"$1.52"}
### Output sample:
[0,0,0,1,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,1,1,1,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,1,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,1,0,0,0,0,0,0,0,1,0,0,0,0,0,0,1,0,1,1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0.6562261013090613,0.044426982996740784]

## 3) redis_loader.py
  This is the python program to load the test case into Redis as workload setup

## 4) docker
  TFS and REST API server images links
