from builder.utils import execer
from builder.utils import log


class Tester(object):

    def __init__(self):
        super(Tester, self).__init__()

    def test(self) -> bool:
        if not self._prepare():
            return False
        if not self._run_test():
            return False
        if not self._gen_report():
            return False
        return True

    def _prepare(self) -> bool:
        ret, _, err = execer.exec('rm -fr unittest')
        if ret != 0:
            log.logger.error(f'clean unittest failed: {err}')
            return False
        ret, _, err = execer.exec('mkdir -p unittest')
        if ret != 0:
            log.logger.error(f'make dir for unittest failed: {err}')
            return False
        return True

    def _run_test(self) -> bool:
        log.logger.info(f'running unittest...')
        ret, _, err = execer.exec('gotestsum --junitfile ./unittest/junit.xml --jsonfile ./unittest/gounit.json\
                  -- -coverpkg=./... -coverprofile=./unittest/cover.out -timeout=3m ./... > ./unittest/gotestsum.txt 2>&1')
        if ret != 0:
            _, out, _ = execer.exec('cat unittest/gotestsum.txt | grep -Ei "fail|ERROR|go:"')
            log.logger.error(f'go test failed: {out}') if ret == 1 else log.logger.error(f'gotestsum failed: {out}')
            return False
        return True

    def _gen_report(self) -> bool:
        ret, _, err = execer.exec('go-test-html-report -f ./unittest/gounit.json -o ./unittest')
        if ret != 0:
            log.logger.error(f'generate overview.html failed: {err}')
            return False
        ret, _, err = execer.exec('mv ./unittest/report.html ./unittest/overview.html')
        if ret != 0:
            log.logger.error(f'rename overview.html failed: {err}')
            return False
        ret, _, err = execer.exec('junit2html ./unittest/junit.xml ./unittest/junit.html')
        if ret != 0:
            log.logger.error(f'generate junit.html failed: {err}')
            return False
        ret, _, err = execer.exec('go tool cover -html=./unittest/cover.out -o ./unittest/cover.html')
        if ret != 0:
            log.logger.error(f'generate cover.html failed: {err}')
            return False
        ret, _, err = execer.exec('go tool cover -func=./unittest/cover.out -o ./unittest/cover.txt')
        if ret != 0:
            log.logger.error(f'generate cover.txt failed: {err}')
            return False
        return True