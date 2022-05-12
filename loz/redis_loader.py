import numpy as np
import json
import numpy
import requests
import joblib
import math
import os
import pandas as pd
import jaydebeapi
import os
#
# Mapper Preparing
def timeEncoder(X):
    X_hm = X['Time'].str.split(':', expand=True)
    d = pd.to_datetime(dict(year=X['Year'],month=X['Month'],day=X['Day'],hour=X_hm[0],minute=X_hm[1])).astype(int)
    return pd.DataFrame(d)

def amtEncoder(X):
    amt = X.apply(lambda x: x[1:]).astype(float).map(lambda amt: max(1,amt)).map(math.log)
    return pd.DataFrame(amt)

def decimalEncoder(X,length=5):
    dnew = pd.DataFrame()
    for i in range(length):
        dnew[i] = np.mod(X,10) 
        X = np.floor_divide(X,10)
    return dnew

def fraudEncoder(X):
    return np.where(X == 'Yes', 1, 0).astype(int)
#
from sklearn_pandas import DataFrameMapper
from sklearn.preprocessing import LabelEncoder
from sklearn.preprocessing import OneHotEncoder
from sklearn.preprocessing import FunctionTransformer
from sklearn.preprocessing import MinMaxScaler
from sklearn.preprocessing import LabelBinarizer
from sklearn.impute import SimpleImputer
#
mapper = joblib.load(open(os.path.join('./','fitted_mapper.pkl'),'rb'))
#
# Reading CSV
#ddf = pd.read_csv('./test_220_100k_os.csv', dtype={"Merchant Name":"str"})
ddf = pd.read_csv('./test_220_100k_os.csv', dtype={"Merchant Name":"str"}, index_col='Index')
indices = np.loadtxt('test_220_100k.indices',dtype=int)
seq_length = 7
#
print(type(ddf),type(indices))#<class 'pandas.core.frame.DataFrame'> <class 'numpy.ndarray'> 
#
def gen_test_batch(ddf, mapper, indices):
    rows = indices.shape[0] 
    for i in range(rows - 1): 
        #print(type(indices[i]),indices[i])
        temp_input = ddf.loc[range(indices[i]-seq_length+1,indices[i]+1)]
        #print("temp_input",temp_input) 
        full_df = mapper.transform(temp_input)
        #print('full_df',full_df)
        tdf = full_df.drop(['Is Fraud?'],axis=1)
        #
        xbatch = tdf.to_numpy().reshape(1, seq_length, -1)
        xbatch_t = np.transpose(xbatch, axes=(1,0,2))
        data = json.dumps({"instances": xbatch_t.tolist()})
        #
        #
        yield indices[i], data
#
import redis
# Change the 9.47.86.127 to any address your Redis Server is
r = redis.StrictRedis(host='9.47.86.127', port=6379)
#
for index,data in gen_test_batch(ddf,mapper,indices):
    print(index,len(data))
    print(r.set(str(index),data))
#
quit()
