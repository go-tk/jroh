from distutils.core import setup

setup(
    name="jroh",
    version="0.0.0",
    description="JSON-RPC over HTTP",
    packages=["jroh"],
    package_dir={"jroh": "src/jroh"},
    package_data={
        "jroh": [
            "data/go/apicommon/errors.go",
            "data/go/apicommon/rpcinfo.go",
            "data/go/apicommon/utils.go",
        ]
    },
    entry_points={
        "console_scripts": [
            "jrohc=jroh.compiler:main",
        ],
    },
)
