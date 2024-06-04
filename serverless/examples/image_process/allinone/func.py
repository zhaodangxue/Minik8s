import couchdb
import cv2
import numpy as np
from uuid import uuid4

def get_image_from_attachment(db, doc_id, attachment_name):
    doc = db[doc_id]
    res = db.get_attachment(doc, attachment_name)
    img = np.asarray(bytearray(res.read()), dtype="uint8")
    img = cv2.imdecode(img, cv2.IMREAD_COLOR)
    return img

def main(params):
    # Connect to CouchDB
    couch = couchdb.Server()
    couch.resource.credentials = (params['COUCHDB_USERNAME'], params['COUCHDB_PASSWORD'])
    src_doc_id = params['id']
    db = couch['images']
    # Get the image from CouchDB
    img = get_image_from_attachment(db, src_doc_id, 'image')

    target_doc = {'_id': uuid4().hex}
    target_doc['metadata'] = {'image_size': img.shape, 'image_dtype': str(img.dtype)}
    target_doc['src_doc'] = src_doc_id
    db.save(target_doc)

    # Process the image
    img = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    # Detect cat face
    cat_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + 'haarcascade_frontalcatface.xml')
    cats = cat_cascade.detectMultiScale(img, scaleFactor=1.1, minNeighbors=5, minSize=(30, 30))
    for (x, y, w, h) in cats:
        img = cv2.rectangle(img, (x, y), (x+w, y+h), (255, 0, 0), 2)
    # Save the processed image
    # Save as png
    img_processed = cv2.imencode('.png', img)[1].tostring()
    cat = {'count': len(cats)}
    target_doc['cat'] = cat
    db.save(target_doc)
    db.put_attachment(target_doc, img_processed, 'image-processed')

    # Create thumbnail
    target_doc = db[target_doc['_id']]
    thumbnail = cv2.resize(img, (100, 100))
    target_doc['thumbnail'] = {'size': thumbnail.shape}
    db.save(target_doc)
    img_thumbnail = cv2.imencode('.png', thumbnail)[1].tostring()
    db.put_attachment(target_doc, img_thumbnail, 'thumbnail')

    return {}

if __name__ == '__main__':
    main({"COUCHDB_USERNAME": "admin", "COUCHDB_PASSWORD": "QWasd123f", "id": "845687d8c2ca4433a4121ae09b4eb441"})
