import os

from builder import base

_REPORT_FILE: str = os.path.join(base.PROJECT_PATH, 'code_check.txt')


class Reporter(object):

    def __init__(self):
        super(Reporter, self).__init__()
        self._report_file = _REPORT_FILE

    def _write_report(self, line: str):
        fi = open(self._report_file, 'a')
        fi.write(line)

    def _get_report(self) -> str:
        return self._report_file
