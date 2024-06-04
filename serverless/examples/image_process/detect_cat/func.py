import couchdb
import cv2
import numpy as np

"""
params:{
    couchdb_username: "",
    couchdb_password: "",
    id: "",
    new_doc_id: "",
    image_size: {
        width: 0,
        height: 0
    }
}
res:{
    couchdb_username: "",
    couchdb_password: "",
    id: "",
    new_doc_id: "",
    image_size: {
        width: 0,
        height: 0
    }
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
    img_gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

    # Detect cat face
    cat_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + 'haarcascade_frontalcatface.xml')
    cats = cat_cascade.detectMultiScale(img_gray, scaleFactor=1.1, minNeighbors=5, minSize=(30, 30))
    for (x, y, w, h) in cats:
        img = cv2.rectangle(img, (x, y), (x+w, y+h), (255, 0, 0), 2)

    # Attach the image to the new document
    target_doc = db[params['new_doc_id']]
    data = cv2.imencode('.png', img)[1].tobytes()
    db.put_attachment(target_doc, data, 'processed_image', 'image/png')

    return params