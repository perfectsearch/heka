import os
import subprocess
import sys

def run_heka():
    heka_process = subprocess.Popen('./ps_build.sh')
    return heka_process.wait()

def main():
    return run_heka()

if __name__ == '__main__':
    sys.exit(main())
