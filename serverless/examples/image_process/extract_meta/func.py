from uuid import uuid4
import couchdb
import cv2
import numpy as np

"""
params:{
    couchdb_username: "",
    couchdb_password: "",
    id: ""
}
res:{
    couchdb_username: "",
    couchdb_password: "",
    id: "",
    new_doc_id: ""
}
"""

def get_image_from_attachment(db, doc_id, attachment_name):
    doc = db[doc_id]
    res = db.get_attachment(doc, attachment_name)
    img = np.asarray(bytearray(res.read()), dtype="uint8")
    img = cv2.imdecode(img, cv2.IMREAD_COLOR)
    return img

def main(params):
    # Connect to CouchDB
    couch = couchdb.Server()
    couch.resource.credentials = (params['couchdb_username'], params['couchdb_password'])
    src_doc_id = params['id']
    db = couch['images']
    
    # Get the image from CouchDB
    img = get_image_from_attachment(db, src_doc_id, 'image')

    # Save the metadata
    target_doc = {'_id': uuid4().hex}
    target_doc['metadata'] = {'image_size': img.shape, 'image_dtype': str(img.dtype)}
    target_doc['src_doc'] = src_doc_id
    db.save(target_doc)

    res = params
    res['new_doc_id'] = target_doc['_id']
    return res