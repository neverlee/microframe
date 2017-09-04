#!/usr/bin/python3

import os
import sys

if __name__ == '__main__':
    exclude = {"cookie115"}
    exclude = {}
    dirs = [adir for adir in os.listdir("./") if os.path.isdir(adir) and (adir not in exclude)]

    header = """package internal

import (
	"github.com/neverlee/microframe/config"
	"github.com/neverlee/microframe/pluginer"
"""
    print(header)
    for adir in dirs:
        print("\t\"github.com/neverlee/microframe/internal/{}\"".format(adir))
    print(")")
    print("")
    print('var Plugins = map[string]func(*config.SrvConf, *config.RawYaml) (pluginer.SrvPluginer, error){')
    for adir in dirs:
        print("\t\"{}\": {}.NewPlugin,".format(adir, adir))
    print("}")