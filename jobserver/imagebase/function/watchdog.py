import os
import json
import job
import threading

from flask import Flask, request

app = Flask(__name__)

user_job = job.Job()
mutex = threading.Lock()

@app.route("/run", methods=['POST'])
def work():
    mutex.acquire()

    try:
        params = json.loads(request.get_data())
        res = user_job.get_status()
        if res['status'] == 'running' or res['status'] == 'success' or res['status'] == 'failed':
            res = {'error': 'job is ' + res['status']}
        else:
            new_thread = threading.Thread(target=user_job.run, args=(params,))
            new_thread.start()
            res = {'status': 'running'}
    except Exception as e:
        res = {'error': str(e)}
    
    mutex.release()
    return json.dumps(res)

@app.route("/status", methods=['GET'])
def status():
    res = user_job.get_status()
    return json.dumps(res)

if __name__ == '__main__':
    app.run(debug=True, host="0.0.0.0", port=int(os.environ.get('PORT', 8080)))
