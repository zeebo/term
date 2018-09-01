package term

type cellbuf struct {
	width  int
	height int
	data   []rune
}

func (cb *cellbuf) init(width, height int) {
	cb.width = width
	cb.height = height
	cb.data = make([]rune, width*height)
}

func (cb *cellbuf) resize(width, height int) {
	if cb.width == width && cb.height == height {
		return
	}

	oldw := cb.width
	oldh := cb.height
	olddata := append([]rune(nil), cb.data...)

	cb.init(width, height)
	cb.clear()

	minw, minh := oldw, oldh

	if width < minw {
		minw = width
	}
	if height < minh {
		minh = height
	}

	for i := 0; i < minh; i++ {
		srco, dsto := i*oldw, i*width
		src := olddata[srco : srco+minw]
		dst := cb.data[dsto : dsto+minw]
		copy(dst, src)
	}
}

func (cb *cellbuf) clear() {
	for i := range cb.data {
		cb.data[i] = ' '
	}
}
