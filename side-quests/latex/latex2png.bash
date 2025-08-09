#!/usr/bin/env bash
FILE_NAME=${1//.latex}

pdflatex $1

convert -density 2000 my.pdf -quality 90 -alpha off -resize 50% file.png
