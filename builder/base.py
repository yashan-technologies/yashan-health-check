import os

from builder.utils import execer
from builder.utils import log

PROJECT_PATH = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))


class BaseBuilder(object):

    def __init__(self, name):
        super(BaseBuilder, self).__init__()
        self.name = name

    def exec_cmd(self, cmd, op):
        log.logger.info("{} {} starting...".format(op, self.name))
        ret, _, err = execer.exec(cmd)
        if ret == 0:
            log.logger.info("{} {} finish!".format(op, self.name))
        else:
            log.logger.error("{} {} failed: {}".format(op, self.name, err))
        return ret
