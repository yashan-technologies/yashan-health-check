import os
import sys

from builder import base, builder
from builder.code_check import checker
from builder.code_test import tester
from builder.utils import log


def clean(args):
    if builder.YHCBuilder().clean() != 0:
        return False


def build(args):
    _update_submodules()
    targets = _gen_targets(args)
    if 'force' in targets:
        if not 'skip_build_template' in targets and builder.YHCBuilder().build_template() != 0:
            return False
        if not 'skip_build_wordgenner' in targets and builder.YHCBuilder().build_wordgenner() != 0:
            return False
        return True if builder.YHCBuilder().force_build() == 0 else False
    if not 'skip_check' in targets and not check(args):
        return False
    if not 'skip_test' in targets and not test(args):
        return False
    if not 'skip_build_template' in targets and builder.YHCBuilder().build_template() != 0:
        return False
    if not 'skip_build_wordgenner' in targets and builder.YHCBuilder().build_wordgenner() != 0:
        return False
    if 'clean' in targets and builder.YHCBuilder().clean() != 0:
        return False
    return True if builder.YHCBuilder().build() == 0 else False


def check(args):
    _format_goimports = False
    targets = _gen_targets(args)
    if "format_goimports" in targets:
        _format_goimports = True
    passed = checker.Checker(_format_goimports).check()
    if not passed:
        log.logger.error('check code failed, please check "code_check.txt" for reason.')
    return passed


def test(args):
    t = tester.Tester()
    result = True
    if not t.test():
        result = False
    log.logger.info('unit test results has been saved to: {}'.format(os.path.join(base.PROJECT_PATH, "unittest")))
    return result


def _update_submodules():
    project_path = base.PROJECT_PATH
    os.chdir(project_path)

    # auto sync the code from git server
    if os.system('git --version') != 0:
        print("command 'git' not found")
        sys.exit(1)

    git_ctrl_path = os.path.join(project_path, ".git")
    if not os.path.exists(git_ctrl_path):
        print("the source code is not under control of git")
        sys.exit(1)

    update_cmd = 'git submodule update --init --recursive'

    ret = os.system(update_cmd)
    if ret != 0:
        sys.exit(1)


def _gen_targets(args):
    running_targets = []
    for arg in dir(args):
        if arg.startswith('_') or arg in ['func', 'jobs']:
            continue
        if not getattr(args, arg, None):
            continue
        running_targets.append(arg)
    return running_targets
