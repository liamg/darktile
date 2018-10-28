package sixel

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/image/bmp"
)

// from https://en.wikipedia.org/wiki/Sixel
func TestParsing(t *testing.T) {

	//raw := `q"1;1;16;16#0;2;0;0;0#1;2;94;75;22#2;2;97;78;31#3;2;97;82;35#4;2;97;82;44#5;2;94;78;25#6;2;91;78;41#7;2;69;60;38#8;2;56;50;35#9;2;63;56;35#10;2;41;38;31#0NB@@!8?@@BN$#1oCA?@!6?@?ACo$#3?O??A?!4@?A??O$#4?_w{{!6}{{w_$#5?G#2CA?@!4?@?AC#5G-#1{_#6K!4?__!4?K#1_{$#5B#4FRrrrz^^zrrrRF#5B$#3?G_#7CCGC??CGCC#3_G$#2?O#9?G!8?G#2?O$#8!4?GC!4?CG-#0NKGG!8?GGKN$#1?BFE!8KEFB$#3???@!8?@$#4!4?@@!4?@@$#5!4?A!6?A$#2!5?A!4?A$#7!6?A??A$#10!6?!4@$#6!7?AA`
	raw := `q
	#0;2;0;0;0#1;2;100;100;0#2;2;0;100;0
	#1~~@@vv@@~~@@~~$
	#2??}}GG}}??}}??-
	#1!14@`
	six, err := ParseString(raw)
	require.Nil(t, err)

	img := six.RGBA()
	require.NotNil(t, img)

	var imageBuf bytes.Buffer
	err = bmp.Encode(io.Writer(&imageBuf), img)

	if err != nil {
		log.Panic(err)
	}

	// Write to file.
	fo, err := os.Create("img.bmp")
	if err != nil {
		panic(err)
	}
	defer fo.Close()
	fo.Write(imageBuf.Bytes())

}
