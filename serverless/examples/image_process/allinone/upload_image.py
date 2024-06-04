import couchdb
import os
import sys
from uuid import uuid4

def main(params):
    # 从环境变量中获取数据库的用户名和密码
    username = os.environ['COUCHDB_USERNAME']
    password = os.environ['COUCHDB_PASSWORD']
    path = params['path']
    # 连接CouchDB
    couch = couchdb.Server()
    couch.resource.credentials = (username, password)
    db = couch['images']
    # 保存文件到CouchDB
    with open(path, 'rb') as f:
        uid = uuid4().hex
        doc = {"_id":uid}
        res = db.save(doc)
        db.put_attachment(doc, f, "image")

    return {}

if __name__ == '__main__':
    # 从参数中获取图片的Path
    path = "/home/varia/tmp/test-img.png"
    main({'path': path})
