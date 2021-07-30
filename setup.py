from distutils.core import setup

setup(
    name="jroh",
    version="0.0.0",
    description="JSON-RPC over HTTP",
    packages=["jroh"],
    package_dir={"jroh": "src/jroh"},
    entry_points={
        "console_scripts": [
            "jrohc=jroh.compiler:main",
        ],
    },
)
