# This chunk is used in generating pylist.json from requirements.txt.
#
# Programatically it will be copied to tmp folder and then following command will be executed:
# cat /tmp/generate_pylist.py | python - requirements.txt pylist.json

import pkg_resources as pr 
import json,sys

gd=pr.get_distribution
res=list()

for i in open(sys.argv[1]):
    try:
        rs={}
        I=gd(i)
        rs["package"]=I.key
        rs["version"]=I.version
        rs["deps"]=set()
        for j in pr.require(i):
            for k in j.requires():
                K=gd(k)
                rs["deps"].add((K.key, K.version))
        rs["deps"] = [{"package":p,"version":v} for p,v in rs["deps"]]
        res.append(rs)
    except: 
        pass
a=sys.argv[2:3]
op=open(a[0],"w")if a else sys.stdout
json.dump(res,op)
op.close()
