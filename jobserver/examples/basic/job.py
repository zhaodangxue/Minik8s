
class Job:
    def __init__(self):
        self.status = 'pending'
    
    def run(self, params):
        self.status = 'success'
        return

    def get_status(self):
        if self.status == 'success' or self.status == 'failed':
            return {'status': self.status, 'result': 'This is a result'}
        return {'status': self.status}
