#!/bin/bash
if [ -z "$1" ]; then
    printf '%s\n' "preprocess_single_output: please provide the mode trip files directory"
    exit 1
fi
DATADIR=""
case $1 in
    */) DATADIR=$1;;
    *) DATADIR="/";DATADIR=$1$DATADIR;;
esac
SUFFIX="*.csv"
FULLPATH=$DATADIR$SUFFIX
NUM=0
for f in $FULLPATH; do
    printf '%s\n' $f
    if [ "$NUM" -ne 0 ]; then
        sed -i 1d "$f"
    fi
    NUM=$((NUM + 1))
done

regex="${DATADIR}*.csv"
newfile="${DATADIR}merged.csv"
cat $regex > $newfile
printf '%s\n' "Finished processing files in ${DATADIR}"    