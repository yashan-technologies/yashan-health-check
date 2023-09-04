import os

from builder import base


class YHCBuilder(base.BaseBuilder):

    def __init__(self):
        super(YHCBuilder, self).__init__('yhc')
        self.current_path = base.PROJECT_PATH

    def build(self):
        os.chdir(self.current_path)
        return self.exec_cmd('make build', 'build')

    def clean(self):
        os.chdir(self.current_path)
        return self.exec_cmd('make clean', 'clean')

    def force_build(self):
        os.chdir(self.current_path)
        return self.exec_cmd('make force', 'force build')
