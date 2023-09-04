import subprocess

encoding: str = 'utf-8'


def exec(cmd: str):
    cmd_list = ['bash', '-c', cmd]
    p = subprocess.Popen(cmd_list, shell=False, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    out, err = p.communicate()
    return p.returncode, out.decode(encoding), err.decode(encoding)
