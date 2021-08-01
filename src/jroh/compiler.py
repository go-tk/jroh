import argparse
import os
import sys

from . import parser, resolver, translator

_ = parser, resolver, translator


def main() -> None:
    arg_parser = argparse.ArgumentParser(prog="jrohc")
    arg_parser.add_argument(
        "files", metavar="FILE", type=str, nargs="+", help="a JROH file to compile"
    )
    arg_parser.add_argument(
        "-o",
        "--out",
        metavar="DIR",
        type=str,
        help="output OpenAPI files to this directory",
        required=True,
    )
    args, file = arg_parser.parse_known_args(sys.argv[1:])
    _compile_files(args.files, args.out)


def _compile_files(file_paths: list[str], output_dir_path: str) -> None:
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
        file_path = os.path.join(output_dir_path, file_path)
        os.makedirs(os.path.dirname(file_path), exist_ok=True)
        with open(file_path, "w+") as f:
            f.write(file_data)


if __name__ == "__main__":
    main()
