import json
import paramiko


class Job:
    def __init__(self):
        self.status = 'pending'
        self.working_dir = '/function'
        self.ssh_identity_file = self.working_dir + 'gpu'
        self.client = None
    
    def run(self, params):
        self.status = 'success'
        password = params['key_file_password']
        self.client = do_login_ssh('pilogin.hpc.sjtu.edu.cn','stu093',self.ssh_identity_file,password)
        # 发送.cu和.slurm文件
        sftp = self.client.open_sftp()
        sftp.put(self.working_dir + '/matrix_add.cu', '/home/stu093/matrix_add.cu')
        sftp.put(self.working_dir + '/matrix_add.slurm', '/home/stu093/matrix_add.slurm')
        # 执行命令，提交任务到计算节点
        stdin, stdout, stderr = self.client.exec_command('sbatch /home/stu093/matrix_add.slurm')
        # 获取任务ID，为第四个字段
        self.job_id = stdout.read().split()[3]
        self.status = 'running'
        return

    def get_status(self):
        if self.status == 'pending':
            return {'status': self.status, 'message': 'Job is pending. Please send run request.'}
        if self.client == None:
            return {'status': self.status, 'message': 'Client is not ready. Please wait for a while.'}
        job_state = self.get_job_status()
        return {'status': self.status, 'job_state': job_state}
    
    def get_job_status(self):
        stdin, stdout, stderr = self.client.exec_command('squeue -u stu093 --json')
        state_json_str = stdout.read()
        state_json = json.loads(state_json_str)
        jobs = state_json['jobs']
        target_job = {}
        for job in jobs:
            if job['job_id'] != self.job_id:
                continue
            target_job['job_id'] = job['job_id']
            target_job['job_state'] = job['job_state']
            target_job['job_output_file'] = job['standard_output']
            # 替换%j为job_id
            target_job['job_output_file'] = target_job['job_output_file'].replace('%j', target_job['job_id'])
        if target_job['job_state'] == 'COMPLETED':
            self.status = 'success'
            # 通过cat读取输出文件内容，存放在output字段中
            stdin, stdout, stderr = self.client.exec_command('cat ' + target_job['job_output_file'])
            target_job['output'] = stdout.read()
        return target_job                       
            
            
def do_login_ssh(hostname, username, ssh_identity_file, key_file_password):
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    private_key = paramiko.Ed25519Key.from_private_key_file(ssh_identity_file, password=key_file_password)
    client.connect(hostname, username=username, pkey=private_key)
    return client

