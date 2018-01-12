#!/bin/bash
if [ -z "$1" ]; then
    printf '%s\n' "preprocess: please provide the mode trip files directory"
    exit 1
fi
DATADIR=""
case $1 in
    */) DATADIR=$1;;
    *) DATADIR="/";DATADIR=$1$DATADIR;;
esac
SUFFIX="*_*.csv"
FULLPATH=$DATADIR$SUFFIX
CUR=""
PREV=""
for f in $FULLPATH; do
    printf '%s\n' $f
    if [[ "$PREV" != "$CUR" && -n "$PREV" ]]; then
        regex="${DATADIR}${PREV}_*.csv"
        newfile="${DATADIR}${PREV}.csv"
        cat $regex > $newfile
        rm $regex
    fi
    PREV=$CUR
    fn=$(basename $f ".csv")
    read PREFIX NUM <<<$(IFS="_"; echo $fn)
    CUR=$PREFIX
    if [ "$NUM" != "1" ]; then
        sed -i 1d "$f"
    fi
    echo "" >> "$f"
done
if [ "$PREV" != "$CUR" ]; then
    regex="${DATADIR}${PREV}_*.csv"
    newfile="${DATADIR}${PREV}.csv"
    cat $regex > $newfile
    rm $regex
    regex="${DATADIR}${CUR}_*.csv"
    newfile="${DATADIR}${CUR}.csv"
    cat $regex > $newfile
    rm $regex
fi
printf '%s\n' "Finished processing files in ${DATADIR}"    
