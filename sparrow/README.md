# Sparrow (jack)

## Desc:

Jack splitted his single-line file and hid the pieces in different places. Can you recover it? 

* File is in format of: "<single-line>\n". 
* You can recover the original file by joining each piece with spaces. (e.g. " ".join(pieces)) 


## Solution:

It's a [steganography](https://en.wikipedia.org/wiki/Steganography) problem:
There is a 1-bit hidden picture in encoded picture such that every pixel of hidden picture, is equal to least significant bit (LSB) of corresponding pixel of encoded picture:

```
Hidden[i, j] = LSB(Encoded[i, j])
```

There are some phrases in every decoded picture; one could use tools like `tesseract` to OCR it. And the flag is the result of concatting and `md5`ing them all.
