rm -f ./wm
cp ../wm/wm ./

[ $? -ne 0 ] && {
    echo "wm not found. $PWD/../wm/wm binary needed"
    exit 1
}

[ ! -e words.txt ] && {
    echo "words.txt not exists"
    exit 1
}

[ ! -e original ] && {
    echo "original/ directory not exists"
    exit 1
}

mkdir -p encoded/
mkdir -p encoded-sorted/
i="0"
for f in $( ls original/* )
do
    echo "Encoding $f"
    i=$(echo "$i + 1" | bc)
    line=$(head words.txt -n $i | tail -n 1)

    rm -f label.bmp
    convert -background black -fill white -pointsize 72 "label:$line" label.bmp

    ./wm encode $f label.bmp encoded-sorted/$i.bmp
done

for f in $( ls encoded-sorted/* | shuf)
do
    echo "converting $f..."
    convert $f encoded/$RANDOM.bmp
done


