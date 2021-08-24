import argparse
import os
import subprocess
import sys
from typing import Optional

from . import go_generator, parser, resolver, translator

_PROG = "jrohc"


def main() -> None:
    arg_parser = argparse.ArgumentParser(prog=_PROG)
    arg_parser.add_argument(
        "files", metavar="FILE", type=str, nargs="+", help="a JROH file to compile"
    )

    def out(out: str) -> str:
        if out == "":
            raise ValueError()
        return out

    arg_parser.add_argument(
        "-o",
        "--out",
        metavar="DIR",
        type=out,
        help="output OpenAPI files to the directory",
        required=True,
    )

    def go_out(go_out: str) -> str:
        i = go_out.find(":")
        if i < 0:
            raise ValueError()
        if i in (0, len(go_out) - 1):
            raise ValueError()
        if go_out.find(":", i + 1) >= 0:
            raise ValueError()
        return go_out

    arg_parser.add_argument(
        "--go_out",
        metavar="DIR:PACKAGE",
        type=go_out,
        help="output go code to the directory as a package",
    )
    args, file = arg_parser.parse_known_args(sys.argv[1:])
    _compile_files(args.files, args.out, args.go_out)


def _compile_files(file_paths: list[str], out: str, go_out: Optional[str]) -> None:
    file_paths.sort()
    file_path_2_file_data: dict[str, str] = {}
    for file_path in file_paths:
        with open(file_path, "r") as f:
            file_data = f.read()
        file_path_2_file_data[file_path] = file_data
        parser.parse_files(file_path_2_file_data)
    results1 = parser.parse_files(file_path_2_file_data)
    for node_uri in results1.ignored_node_uris:
        print(f"WARNING: node ignored: {node_uri}", file=sys.stderr)
    results2 = resolver.resolve_specs(results1.specs)
    for node_uri in results2.unused_node_uris:
        print(f"WARNING: node unused: {node_uri}", file=sys.stderr)
    results3 = translator.translate_specs(results2.merged_specs)
    for file_path, file_data in results3.file_path_2_file_data.items():
        output_dir_path = out
        file_path = os.path.join(output_dir_path, file_path)
        os.makedirs(os.path.dirname(file_path), exist_ok=True)
        with open(file_path, "w+") as f:
            f.write(f"# File generated by {_PROG}. DO NOT EDIT.\n")
            f.write(file_data)
    if go_out is not None:
        output_dir_path, output_package_path = go_out.split(":", 1)
        results4 = go_generator.generate_code(
            output_package_path, results2.merged_specs
        )
        file_paths = []
        for file_path, file_data in results4.file_path_2_file_data.items():
            file_path = os.path.join(output_dir_path, file_path)
            file_paths.append(file_path)
            os.makedirs(os.path.dirname(file_path), exist_ok=True)
            with open(file_path, "w+") as f:
                f.write(f"// Code generated by {_PROG}. DO NOT EDIT.\n\n")
                f.write(file_data)
        _format_go_code(file_paths)


def _format_go_code(file_paths: list[str]) -> None:
    try:
        subprocess.run(["gofmt", "-w", *file_paths])
    except Exception as e:
        print(f"WARNING: go code formatting failed: {e}", file=sys.stderr)


if __name__ == "__main__":
    main()
