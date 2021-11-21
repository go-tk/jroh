from distutils.core import setup

setup(
    name="jroh",
    version="0.7.0",
    description="JSON-RPC over HTTP",
    packages=["jroh"],
    package_dir={"jroh": "src/jroh"},
    python_requires=">=3.9",
    install_requires=[
        "PyYAML==5.4.1",
        "Mako==1.1.4",
        "google-re2==0.2.20210901",
    ],
    entry_points={
        "console_scripts": [
            "jrohc=jroh.compiler:main",
        ],
    },
)
