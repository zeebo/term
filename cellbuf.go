package term

type cellbuf struct {
	width  int
	height int
	data   []rune
}

func (this *cellbuf) init(width, height int) {
	this.width = width
	this.height = height
	this.data = make([]rune, width*height)
}

func (this *cellbuf) resize(width, height int) {
	if this.width == width && this.height == height {
		return
	}

	oldw := this.width
	oldh := this.height
	olddata := append([]rune(nil), this.data...)

	this.init(width, height)
	this.clear()

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
		dst := this.data[dsto : dsto+minw]
		copy(dst, src)
	}
}

func (this *cellbuf) clear() {
	for i := range this.data {
		this.data[i] = ' '
	}
}
