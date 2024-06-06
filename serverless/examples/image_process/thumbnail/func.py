import couchdb
import cv2
import numpy as np

"""
params:{
    couchdb_host: "",
    couchdb_username: "",
    couchdb_password: "",
    id: "",
    new_doc_id: ""
}
res:{
    res_doc_id: ""
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
    couch = couchdb.Server(params['couchdb_host'])
    couch.resource.credentials = (params['couchdb_username'], params['couchdb_password'])
    target_doc_id = params['new_doc_id']
    db = couch['images']

    # Get the image from CouchDB
    img = get_image_from_attachment(db, target_doc_id, 'processed_image')

    # Generate thumbnail
    # Scale till the smaller side is 100
    scale = 100 / min(img.shape[:2])
    thumbnail = cv2.resize(img, (0, 0), fx=scale, fy=scale)
    target_doc = db[target_doc_id]
    target_doc['thumbnail'] = {'size': thumbnail.shape}
    db.save(target_doc)
    data = cv2.imencode('.png', thumbnail)[1].tobytes()
    db.put_attachment(target_doc, data, 'thumbnail', 'image/png')

    return {'res_doc_id': target_doc_id}
