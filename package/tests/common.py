from dataclasses import dataclass, field
from typing import Optional, Type
from unittest import TestCase

from ..jroh import dumper, parser, resolver


@dataclass
class TestData:
    in_file_path_2_file_data: dict[str, str] = field(default_factory=dict)

    out_exception_type: Optional[Type[BaseException]] = None
    out_exception_str: str = ""
    out_ignored_node_uris: Optional[set[str]] = None
    out_unused_node_uris: Optional[set[str]] = None
    out_file_path_2_file_data: Optional[dict[str, str]] = None


def test(test_case: TestCase, test_data_list: list[TestData]) -> None:
    for i, test_data in enumerate(test_data_list):
        with test_case.subTest(i):
            if test_data.out_exception_type is None:
                result1 = parser.parse_files(test_data.in_file_path_2_file_data)
                result2 = resolver.resolve_specs(result1.specs)
                result3 = dumper.dump_specs(result2.merged_specs)
                if (ignored_node_uris := test_data.out_ignored_node_uris) is not None:
                    test_case.assertSetEqual(
                        set(result1.ignored_node_uris), ignored_node_uris
                    )
                if (unused_node_uris := test_data.out_unused_node_uris) is not None:
                    test_case.assertSetEqual(
                        set(result2.unused_node_uris), unused_node_uris
                    )
                if (
                    file_path_2_file_data := test_data.out_file_path_2_file_data
                ) is not None:
                    test_case.assertDictEqual(
                        result3.file_path_2_file_data, file_path_2_file_data
                    )
            else:
                with test_case.assertRaisesRegex(
                    test_data.out_exception_type, test_data.out_exception_str
                ):
                    result1 = parser.parse_files(test_data.in_file_path_2_file_data)
                    result2 = resolver.resolve_specs(result1.specs)
                    result3 = dumper.dump_specs(result2.merged_specs)
                    _ = result3
