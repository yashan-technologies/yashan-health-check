#!/usr/bin/python
# -*- coding: UTF-8 -*-
import argparse
import sys

from builder import handler


class HelpFormatter(argparse.ArgumentDefaultsHelpFormatter):

    def __init__(self, prog, indent_increment=1, max_help_position=48, width=256):
        super(HelpFormatter, self).__init__(prog,
                                            indent_increment=indent_increment,
                                            max_help_position=max_help_position,
                                            width=width)


def get_parser():
    parser = argparse.ArgumentParser(
        description="yhc building tool",
        prog="build.py",
        usage='%(prog)s [-h, --help]',
        formatter_class=argparse.HelpFormatter,
    )
    subparser = parser.add_subparsers()
    set_clean_argument(subparser)
    set_build_argument(subparser)
    set_check_argument(subparser)
    set_test_argument(subparser)
    return parser


def set_build_argument(subparser):
    sp = subparser.add_parser(
        "build", help="build yhc", formatter_class=HelpFormatter)
    sp.add_argument("--skip-check", action="store_true",
                    default=False, help="build without checking code")
    sp.add_argument("--skip-test", action="store_true",
                    default=False, help="build without running unit test")
    sp.add_argument("--format-goimports", action="store_true",
                    default=False, help="format go imports")
    sp.add_argument("-c", "--clean", action="store_true",
                    default=False, help="clean before building")
    sp.add_argument("--skip-build-template", action="store_true",
                    default=False, help="skip build html template before build")
    sp.add_argument("--skip-build-wordgenner", action="store_true",
                    default=False, help="skip build word genner before build")
    sp.add_argument("-f",
                    "--force",
                    action="store_true",
                    default=False,
                    help="clean before building, then build without checking code and running unit test")
    sp.set_defaults(func=build)


def set_clean_argument(subparser):
    sp = subparser.add_parser(
        "clean", help="clean yhc", formatter_class=HelpFormatter)
    sp.set_defaults(func=clean)


def set_check_argument(subparser):
    sp = subparser.add_parser(
        "check", help="check code", formatter_class=HelpFormatter)
    sp.add_argument("--format-goimports", action="store_true",
                    default=False, help="format go imports")
    sp.set_defaults(func=check)


def set_test_argument(subparser):
    sp = subparser.add_parser(
        "test", help="run unit test", formatter_class=HelpFormatter)
    sp.set_defaults(func=test)


def build(args):
    if not handler.build(args):
        return 1
    return 0


def clean(args):
    if not handler.clean(args):
        return 1
    return 0


def check(args):
    if not handler.check(args):
        print('Check code failed, please check "code_check.txt" for reason.')
        return 1
    return 0


def test(args):
    if not handler.test(args):
        return 1
    return 0


if __name__ == "__main__":
    parser = get_parser()
    args = parser.parse_args()
    if hasattr(args, 'func') and args.func is not None:
        sys.exit(args.func(args))
    parser.print_usage()
    sys.exit(0)
