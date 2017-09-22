#! /bin/bash
set -e

scp -r dist/* zero:www \
    && ssh zero 'mv ~/www/door{,.old} && mv ~/www/door{.new,} && service door restart'
