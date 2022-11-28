#!/bin/env python3

import os
import subprocess
import yaml


def main():
    with open("toto.yaml") as yaml_file:
        yaml_content = yaml.load(yaml_file)
        subprocess.run(yaml_content["toto"], check=False)


if __name__ == "__main__":
    main()
